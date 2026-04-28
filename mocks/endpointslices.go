package mocks

import (
	"context"
	"fmt"

	v1 "k8s.io/api/discovery/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	applyconfigurationsv1 "k8s.io/client-go/applyconfigurations/discovery/v1"
	discoveryv1 "k8s.io/client-go/kubernetes/typed/discovery/v1"
)

type esClient struct{}

const (
	MockEndpointSliceError       = "mock endpointslice didnt work"
	NotReadyEndpointsServiceName = "not-ready-endpoints"
)

func (es esClient) Create(
	ctx context.Context,
	endpointSlice *v1.EndpointSlice,
	opts metav1.CreateOptions,
) (*v1.EndpointSlice, error) {
	return nil, fmt.Errorf("not implemented")
}

func (es esClient) Update(
	ctx context.Context,
	endpointSlice *v1.EndpointSlice,
	opts metav1.UpdateOptions,
) (*v1.EndpointSlice, error) {
	return nil, fmt.Errorf("not implemented")
}

func (es esClient) Delete(
	ctx context.Context,
	name string,
	opts metav1.DeleteOptions,
) error {
	return fmt.Errorf("not implemented")
}

func (es esClient) DeleteCollection(
	ctx context.Context,
	opts metav1.DeleteOptions,
	listOpts metav1.ListOptions,
) error {
	return fmt.Errorf("not implemented")
}

func (es esClient) Get(
	ctx context.Context,
	name string,
	opts metav1.GetOptions,
) (*v1.EndpointSlice, error) {
	return nil, fmt.Errorf("not implemented")
}

func (es esClient) List(
	ctx context.Context,
	opts metav1.ListOptions,
) (*v1.EndpointSliceList, error) {
	var endpoints []v1.Endpoint

	switch opts.LabelSelector {
	case fmt.Sprintf("%s=%s", v1.LabelServiceName, SucceedingServiceName):
		ready := true
		endpoints = []v1.Endpoint{
			{
				Addresses:  []string{"127.0.0.1"},
				Conditions: v1.EndpointConditions{Ready: &ready},
			},
		}
	case fmt.Sprintf("%s=%s", v1.LabelServiceName, EmptyEndpointsServiceName):
		endpoints = []v1.Endpoint{}
	case fmt.Sprintf("%s=%s", v1.LabelServiceName, NotReadyEndpointsServiceName):
		notReady := false
		endpoints = []v1.Endpoint{
			{
				Addresses:  []string{"127.0.0.1"},
				Conditions: v1.EndpointConditions{Ready: &notReady},
			},
		}
	default:
		return nil, fmt.Errorf(MockEndpointSliceError)
	}

	endpointslices := []v1.EndpointSlice{
		{
			ObjectMeta: metav1.ObjectMeta{Name: "endpointslice-1"},
			Endpoints:  endpoints,
		},
	}

	return &v1.EndpointSliceList{
		Items: endpointslices,
	}, nil
}

func (es esClient) Watch(
	ctx context.Context,
	opts metav1.ListOptions,
) (watch.Interface, error) {
	return nil, fmt.Errorf("not implemented")
}

func (es esClient) Patch(
	ctx context.Context,
	name string,
	pt types.PatchType,
	data []byte,
	opts metav1.PatchOptions,
	subresources ...string,
) (result *v1.EndpointSlice, err error) {
	return nil, fmt.Errorf("not implemented")
}

func (es esClient) Apply(
	ctx context.Context,
	endpointslices *applyconfigurationsv1.EndpointSliceApplyConfiguration,
	opts metav1.ApplyOptions,
) (result *v1.EndpointSlice, err error) {
	return nil, fmt.Errorf("not implemented")
}

func NewESClient() discoveryv1.EndpointSliceInterface {
	return esClient{}
}
