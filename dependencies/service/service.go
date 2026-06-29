package service

import (
	"context"
	"fmt"

	discoveryv1 "k8s.io/api/discovery/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	entry "opendev.org/airship/kubernetes-entrypoint/entrypoint"
	"opendev.org/airship/kubernetes-entrypoint/util/env"
)

const FailingStatusFormat = "Service %v has no ready endpoints"

type Service struct {
	name      string
	namespace string
}

func init() {
	serviceEnv := fmt.Sprintf("%sSERVICE", entry.DependencyPrefix)
	if serviceDeps := env.SplitEnvToDeps(serviceEnv); serviceDeps != nil {
		if len(serviceDeps) > 0 {
			for _, dep := range serviceDeps {
				entry.Register(NewService(dep.Name, dep.Namespace))
			}
		}
	}
}

func NewService(name string, namespace string) Service {
	return Service{
		name:      name,
		namespace: namespace,
	}
}

func (s Service) IsResolved(ctx context.Context, entrypoint entry.EntrypointInterface) (bool, error) {
	// Try EndpointSlices first (new API; requires endpointslices RBAC in the chart).
	listOptions := metav1.ListOptions{LabelSelector: fmt.Sprintf("%s=%s", discoveryv1.LabelServiceName, s.name)}
	if esList, err := entrypoint.Client().EndpointSlices(s.namespace).List(ctx, listOptions); err == nil {
		for _, es := range esList.Items {
			for _, e := range es.Endpoints {
				if len(e.Addresses) > 0 && (e.Conditions.Ready == nil || *e.Conditions.Ready) {
					return true, nil
				}
			}
		}
	}

	// Fall back to Endpoints (legacy API; required by charts that only grant endpoints RBAC).
	if ep, err := entrypoint.Client().Endpoints(s.namespace).Get(ctx, s.name, metav1.GetOptions{}); err == nil {
		for _, subset := range ep.Subsets {
			if len(subset.Addresses) > 0 {
				return true, nil
			}
		}
	}

	return false, fmt.Errorf(FailingStatusFormat, s.name)
}

func (s Service) String() string {
	return fmt.Sprintf("Service %s in namespace %s", s.name, s.namespace)
}
