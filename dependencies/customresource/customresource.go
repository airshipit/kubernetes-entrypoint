package customresource

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"opendev.org/airship/kubernetes-entrypoint/entrypoint"
	"opendev.org/airship/kubernetes-entrypoint/logger"
	"opendev.org/airship/kubernetes-entrypoint/util/env"
)

// A Resolver represents the state of a CustomResource
type Resolver struct {
	APIVersion string  `json:"apiVersion"`
	Kind       string  `json:"kind"`
	Name       string  `json:"name"`
	Namespace  string  `json:"namespace"`
	Fields     []Field `json:"fields"`
}

var _ entrypoint.Resolver = Resolver{}

// A Field represents a key-value pair
type Field struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func init() {
	crEnv := fmt.Sprintf("%sCUSTOM_RESOURCE", entrypoint.DependencyPrefix)
	resolvers, err := fromEnv(crEnv)
	if err != nil {
		logger.Error.Printf(err.Error())
	}
	for _, resolver := range resolvers {
		entrypoint.Register(resolver)
	}
}

// IsResolved will return true when the values for each key in r.Fields is the same as the resource in the cluster
func (r Resolver) IsResolved(ctx context.Context, ep entrypoint.EntrypointInterface) (bool, error) {
	customResource, err := ep.Client().CustomResource(ctx, r.APIVersion, r.Kind, r.Namespace, r.Name)
	if err != nil {
		return false, err
	}

	for _, field := range r.Fields {
		key := field.Key
		expected := field.Value

		// Extract the specified value from the resource
		actual, found, err := unstructured.NestedFieldNoCopy(customResource.Object, strings.Split(key, ".")...)
		if err != nil {
			return false, err
		}
		if !found {
			return false, fmt.Errorf("could not find key [%s]", key)
		}
		if actual != expected {
			return false, fmt.Errorf("expected value of [%s] to be [%s], but got [%s]", key, expected, actual)
		}
	}

	return true, nil
}

// fromEnv reads the value of the jsonEnv variable and returns the array of
// Resolvers it contains, if any
func fromEnv(jsonEnv string) ([]Resolver, error) {
	resolvers := []Resolver{}
	jsonEnvVal, isSet := os.LookupEnv(jsonEnv)
	if !isSet {
		return resolvers, nil
	}

	err := json.Unmarshal([]byte(jsonEnvVal), &resolvers)
	if err != nil {
		return resolvers, fmt.Errorf("unable to unmarshal variable %s with value %s: %s",
			jsonEnv, jsonEnvVal, err.Error())
	}

	namespace := env.GetBaseNamespace()
	for i := range resolvers {
		if resolvers[i].Namespace == "" {
			resolvers[i].Namespace = namespace
		}
	}

	return resolvers, nil
}
