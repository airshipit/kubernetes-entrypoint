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

// DefaultConditionStatus is the status a Condition expects when it does not
// say otherwise, matching the common case of waiting for something to become
// ready rather than to fail.
const DefaultConditionStatus = "True"

// A Resolver represents the state of a CustomResource
type Resolver struct {
	APIVersion string      `json:"apiVersion"`
	Kind       string      `json:"kind"`
	Name       string      `json:"name"`
	Namespace  string      `json:"namespace"`
	Fields     []Field     `json:"fields"`
	Conditions []Condition `json:"conditions"`
}

var _ entrypoint.Resolver = Resolver{}

// A Field represents a key-value pair
type Field struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// A Condition names an entry of the resource's status.conditions list and the
// status that entry must report, the way kubectl wait --for=condition=<type>
// does.
//
// Conditions cannot be expressed as a Field, because a Field addresses one
// nested key by name and a condition list is addressed by matching the type of
// its elements. Resources that follow the Kubernetes API conventions report
// readiness this way, so without this a dependency on such a resource cannot
// be written at all.
type Condition struct {
	Type string `json:"type"`
	// Status defaults to DefaultConditionStatus when empty. It is defaulted
	// where it is compared rather than where it is parsed, so that a Resolver
	// built in code behaves the same as one read from the environment.
	Status string `json:"status"`
}

func init() {
	crEnv := fmt.Sprintf("%sCUSTOM_RESOURCE", entrypoint.DependencyPrefix)
	resolvers, err := fromEnv(crEnv)
	if err != nil {
		logger.Error.Printf("Error initializing custom resource: %s", err.Error()) // Fixed format string
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

	for _, condition := range r.Conditions {
		expected := condition.Status
		if expected == "" {
			expected = DefaultConditionStatus
		}

		actual, err := conditionStatus(customResource.Object, condition.Type)
		if err != nil {
			return false, err
		}
		if actual != expected {
			return false, fmt.Errorf("expected condition [%s] to be [%s], but got [%s]",
				condition.Type, expected, actual)
		}
	}

	return true, nil
}

// conditionStatus returns the status reported by the condition of the given
// type, or an error when the resource carries no such condition yet. A
// resource that has not been reconciled has no conditions at all, which is not
// a permanent failure: the caller retries.
func conditionStatus(object map[string]interface{}, conditionType string) (string, error) {
	value, found, err := unstructured.NestedFieldNoCopy(object, "status", "conditions")
	if err != nil {
		return "", err
	}
	if !found {
		return "", fmt.Errorf("resource has no [status.conditions]")
	}

	conditions, ok := value.([]interface{})
	if !ok {
		return "", fmt.Errorf("[status.conditions] is not a list")
	}

	for _, item := range conditions {
		condition, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		if condition["type"] != conditionType {
			continue
		}
		status, ok := condition["status"].(string)
		if !ok {
			return "", fmt.Errorf("condition [%s] has no string status", conditionType)
		}
		return status, nil
	}

	return "", fmt.Errorf("could not find condition of type [%s]", conditionType)
}

// fromEnv reads the value of the jsonEnv variable and returns the array of
// Resolvers it contains, if any
func fromEnv(jsonEnv string) ([]Resolver, error) {
	resolvers := []Resolver{}
	jsonEnvVal, isSet := os.LookupEnv(jsonEnv)
	if !isSet {
		return resolvers, nil
	}

	// Check if the environment variable is empty
	if strings.TrimSpace(jsonEnvVal) == "" {
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
