package mocks

import (
	appsv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	v1batch "k8s.io/client-go/kubernetes/typed/batch/v1"
	v1core "k8s.io/client-go/kubernetes/typed/core/v1"

	cli "opendev.org/airship/kubernetes-entrypoint/client"
)

type Client struct {
	v1core.PodInterface
	v1core.ServiceInterface
	appsv1.DaemonSetInterface
	v1core.EndpointsInterface
	v1batch.JobInterface
}

func (c Client) Pods(namespace string) v1core.PodInterface {
	return c.PodInterface
}

func (c Client) Services(namespace string) v1core.ServiceInterface {
	return c.ServiceInterface
}

func (c Client) DaemonSets(namespace string) appsv1.DaemonSetInterface {
	return c.DaemonSetInterface
}

func (c Client) Endpoints(namespace string) v1core.EndpointsInterface {
	return c.EndpointsInterface
}
func (c Client) Jobs(namespace string) v1batch.JobInterface {
	return c.JobInterface
}

func NewClient() cli.ClientInterface {
	return Client{
		NewPClient(),
		NewSClient(),
		NewDSClient(),
		NewEClient(),
		NewJClient(),
	}
}
