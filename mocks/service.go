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
	restclient "k8s.io/client-go/rest"
)

type sClient struct {
}

const (
	MockServiceError        = "mock service didnt work"
	SucceedingServiceName   = "succeed"
	EmptySubsetsServiceName = "empty-subsets"
	FailingServiceName      = "fail"
)

func (s sClient) Create(ctx context.Context, service *v1.Service, opts metav1.CreateOptions) (*v1.Service, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s sClient) Update(ctx context.Context, service *v1.Service, opts metav1.UpdateOptions) (*v1.Service, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s sClient) UpdateStatus(ctx context.Context, service *v1.Service, opts metav1.UpdateOptions) (*v1.Service, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s sClient) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return fmt.Errorf("not implemented")
}

func (s sClient) Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.Service, error) {
	if name == FailingServiceName {
		return nil, fmt.Errorf(MockServiceError)
	}
	return &v1.Service{
		ObjectMeta: metav1.ObjectMeta{Name: name},
	}, nil
}

func (s sClient) List(ctx context.Context, opts metav1.ListOptions) (*v1.ServiceList, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s sClient) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s sClient) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.Service, err error) {
	return nil, fmt.Errorf("not implemented")
}

func (s sClient) Apply(ctx context.Context, service *corev1applyconfigurations.ServiceApplyConfiguration, opts metav1.ApplyOptions) (result *v1.Service, err error) {
	return nil, fmt.Errorf("not implemented")
}

func (s sClient) ApplyStatus(ctx context.Context, service *corev1applyconfigurations.ServiceApplyConfiguration, opts metav1.ApplyOptions) (result *v1.Service, err error) {
	return nil, fmt.Errorf("not implemented")
}

func (s sClient) ProxyGet(scheme, name, port, path string, params map[string]string) restclient.ResponseWrapper {
	return nil
}

func NewSClient() v1core.ServiceInterface {
	return sClient{}
}
