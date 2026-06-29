package mocks

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	corev1applyconfigurations "k8s.io/client-go/applyconfigurations/core/v1"
	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
)

type epClient struct{}

func (e epClient) Create(
	ctx context.Context,
	endpoints *v1.Endpoints,
	opts metav1.CreateOptions,
) (*v1.Endpoints, error) {
	return nil, fmt.Errorf("not implemented")
}

func (e epClient) Update(
	ctx context.Context,
	endpoints *v1.Endpoints,
	opts metav1.UpdateOptions,
) (*v1.Endpoints, error) {
	return nil, fmt.Errorf("not implemented")
}

func (e epClient) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return fmt.Errorf("not implemented")
}

func (e epClient) DeleteCollection(
	ctx context.Context,
	opts metav1.DeleteOptions,
	listOpts metav1.ListOptions,
) error {
	return fmt.Errorf("not implemented")
}

func (e epClient) Get(
	ctx context.Context,
	name string,
	opts metav1.GetOptions,
) (*v1.Endpoints, error) {
	switch name {
	case SucceedingServiceName:
		return &v1.Endpoints{
			ObjectMeta: metav1.ObjectMeta{Name: name},
			Subsets: []v1.EndpointSubset{
				{Addresses: []v1.EndpointAddress{{IP: "127.0.0.1"}}},
			},
		}, nil
	case EmptyEndpointsServiceName:
		return &v1.Endpoints{
			ObjectMeta: metav1.ObjectMeta{Name: name},
		}, nil
	case NotReadyEndpointsServiceName:
		return &v1.Endpoints{
			ObjectMeta: metav1.ObjectMeta{Name: name},
			Subsets: []v1.EndpointSubset{
				{NotReadyAddresses: []v1.EndpointAddress{{IP: "127.0.0.1"}}},
			},
		}, nil
	default:
		return nil, fmt.Errorf("mock endpoints didnt work")
	}
}

func (e epClient) List(
	ctx context.Context,
	opts metav1.ListOptions,
) (*v1.EndpointsList, error) {
	return nil, fmt.Errorf("not implemented")
}

func (e epClient) Watch(
	ctx context.Context,
	opts metav1.ListOptions,
) (watch.Interface, error) {
	return nil, fmt.Errorf("not implemented")
}

func (e epClient) Patch(
	ctx context.Context,
	name string,
	pt types.PatchType,
	data []byte,
	opts metav1.PatchOptions,
	subresources ...string,
) (*v1.Endpoints, error) {
	return nil, fmt.Errorf("not implemented")
}

func (e epClient) Apply(
	ctx context.Context,
	endpoints *corev1applyconfigurations.EndpointsApplyConfiguration,
	opts metav1.ApplyOptions,
) (*v1.Endpoints, error) {
	return nil, fmt.Errorf("not implemented")
}

func NewEPClient() v1core.EndpointsInterface {
	return epClient{}
}
