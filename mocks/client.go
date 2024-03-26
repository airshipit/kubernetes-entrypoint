package mocks

import (
	"context"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	v1apps "k8s.io/client-go/kubernetes/typed/apps/v1"
	v1batch "k8s.io/client-go/kubernetes/typed/batch/v1"
	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
)

type Client struct {
	v1core.PodInterface
	v1core.ServiceInterface
	v1apps.DaemonSetInterface
	v1core.EndpointsInterface
	v1batch.JobInterface

	FakeCustomResource *unstructured.Unstructured
	Err                error
}

func (c Client) Pods(namespace string) v1core.PodInterface {
	return c.PodInterface
}

func (c Client) Services(namespace string) v1core.ServiceInterface {
	return c.ServiceInterface
}

func (c Client) DaemonSets(namespace string) v1apps.DaemonSetInterface {
	return c.DaemonSetInterface
}

func (c Client) Endpoints(namespace string) v1core.EndpointsInterface {
	return c.EndpointsInterface
}

func (c Client) Jobs(namespace string) v1batch.JobInterface {
	return c.JobInterface
}

func (c Client) CustomResource(
	ctx context.Context,
	apiVersion, namespace, resource, name string,
) (*unstructured.Unstructured, error) {
	return c.FakeCustomResource, c.Err
}

func NewClient() *Client {
	return &Client{
		PodInterface:       NewPClient(),
		ServiceInterface:   NewSClient(),
		DaemonSetInterface: NewDSClient(),
		EndpointsInterface: NewEClient(),
		JobInterface:       NewJClient(),
	}
}
