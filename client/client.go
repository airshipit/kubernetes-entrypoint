package client

import (
	"k8s.io/client-go/kubernetes"
	appsv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	v1batch "k8s.io/client-go/kubernetes/typed/batch/v1"
	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
)

type ClientInterface interface {
	Pods(string) v1core.PodInterface
	Jobs(string) v1batch.JobInterface
	Endpoints(string) v1core.EndpointsInterface
	DaemonSets(string) appsv1.DaemonSetInterface
	Services(string) v1core.ServiceInterface
}
type Client struct {
	*kubernetes.Clientset
}

func (c Client) Pods(namespace string) v1core.PodInterface {
	return c.Clientset.CoreV1().Pods(namespace)
}

func (c Client) Jobs(namespace string) v1batch.JobInterface {
	return c.Clientset.BatchV1().Jobs(namespace)
}

func (c Client) Endpoints(namespace string) v1core.EndpointsInterface {
	return c.Clientset.CoreV1().Endpoints(namespace)
}
func (c Client) DaemonSets(namespace string) appsv1.DaemonSetInterface {
	return c.Clientset.AppsV1().DaemonSets(namespace)
}

func (c Client) Services(namespace string) v1core.ServiceInterface {
	return c.Clientset.CoreV1().Services(namespace)
}

func New(config *rest.Config) (ClientInterface, error) {
	if config == nil {
		var err error
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			return nil, err
		}
		return Client{Clientset: clientset}, nil
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return Client{Clientset: clientset}, nil

}
