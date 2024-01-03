package mocks

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	policyv1 "k8s.io/api/policy/v1"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	corev1applyconfigurations "k8s.io/client-go/applyconfigurations/core/v1"
	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
	restclient "k8s.io/client-go/rest"
)

const MockContainerName = "TEST_CONTAINER"

type pClient struct {
}

const (
	PodNotPresent                   = "NOT_PRESENT"
	PodEnvVariableValue             = "podlist"
	FailingMatchLabel               = "INCORRECT"
	SameHostNotReadyMatchLabel      = "SAME_HOST_NOT_READY"
	SameHostReadyMatchLabel         = "SAME_HOST_READY"
	SameHostSomeReadyMatchLabel     = "SAME_HOST_SOME_READY"
	DifferentHostReadyMatchLabel    = "DIFFERENT_HOST_READY"
	DifferentHostNotReadyMatchLabel = "DIFFERENT_HOST_NOT_READY"
	NoPodsMatchLabel                = "NO_PODS"
)

func (p pClient) Create(ctx context.Context, pod *v1.Pod, opts metav1.CreateOptions) (*v1.Pod, error) {
	return nil, fmt.Errorf("not implemented")
}

func (p pClient) Update(ctx context.Context, pod *v1.Pod, opts metav1.UpdateOptions) (*v1.Pod, error) {
	return nil, fmt.Errorf("not implemented")
}

func (p pClient) UpdateStatus(ctx context.Context, pod *v1.Pod, opts metav1.UpdateOptions) (*v1.Pod, error) {
	return nil, fmt.Errorf("not implemented")
}

func (p pClient) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return fmt.Errorf("not implemented")
}

func (p pClient) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	return fmt.Errorf("not implemented")
}

func (p pClient) Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.Pod, error) {
	if name == PodNotPresent {
		return nil, fmt.Errorf("could not get pod with the name %s", name)
	}

	return &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: name},
		Status: v1.PodStatus{
			ContainerStatuses: []v1.ContainerStatus{
				{
					Name:  MockContainerName,
					Ready: true,
				},
			},
			HostIP: "127.0.0.1",
		},
	}, nil
}

func (p pClient) List(ctx context.Context, opts metav1.ListOptions) (*v1.PodList, error) {
	if opts.LabelSelector == fmt.Sprintf("name=%s", FailingMatchLabel) {
		return nil, fmt.Errorf("Client received incorrect pod label names")
	}

	readyPodSameHost := NewPod(true, "127.0.0.1")
	notReadyPodSameHost := NewPod(false, "127.0.0.1")
	readyPodDifferentHost := NewPod(true, "10.0.0.1")
	notReadyPodDifferentHost := NewPod(false, "10.0.0.1")

	var pods []v1.Pod

	if opts.LabelSelector == fmt.Sprintf("name=%s", SameHostNotReadyMatchLabel) {
		pods = []v1.Pod{notReadyPodSameHost}
	}
	if opts.LabelSelector == fmt.Sprintf("name=%s", SameHostReadyMatchLabel) {
		pods = []v1.Pod{readyPodSameHost, notReadyPodDifferentHost}
	}
	if opts.LabelSelector == fmt.Sprintf("name=%s", SameHostSomeReadyMatchLabel) {
		pods = []v1.Pod{readyPodSameHost, notReadyPodSameHost}
	}
	if opts.LabelSelector == fmt.Sprintf("name=%s", DifferentHostReadyMatchLabel) {
		pods = []v1.Pod{notReadyPodSameHost, readyPodDifferentHost}
	}
	if opts.LabelSelector == fmt.Sprintf("name=%s", DifferentHostNotReadyMatchLabel) {
		pods = []v1.Pod{notReadyPodDifferentHost}
	}
	if opts.LabelSelector == fmt.Sprintf("name=%s", NoPodsMatchLabel) {
		pods = []v1.Pod{}
	}

	return &v1.PodList{
		Items: pods,
	}, nil
}

func (p pClient) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return nil, fmt.Errorf("not implemented")
}

func (p pClient) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.Pod, err error) {
	return nil, fmt.Errorf("not implemented")
}

func (p pClient) Apply(ctx context.Context, pod *corev1applyconfigurations.PodApplyConfiguration, opts metav1.ApplyOptions) (result *v1.Pod, err error) {
	return nil, fmt.Errorf("not implemented")
}

func (p pClient) ApplyStatus(ctx context.Context, pod *corev1applyconfigurations.PodApplyConfiguration, opts metav1.ApplyOptions) (result *v1.Pod, err error) {
	return nil, fmt.Errorf("not implemented")
}

func (p pClient) UpdateEphemeralContainers(ctx context.Context, podName string, pod *v1.Pod, opts metav1.UpdateOptions) (*v1.Pod, error) {
	return nil, fmt.Errorf("not implemented")
}

func (p pClient) Bind(ctx context.Context, binding *v1.Binding, opts metav1.CreateOptions) error {
	return fmt.Errorf("not implemented")
}

func (p pClient) Evict(ctx context.Context, eviction *policyv1beta1.Eviction) error {
	return fmt.Errorf("not implemented")
}

func (p pClient) EvictV1(ctx context.Context, eviction *policyv1.Eviction) error {
	return fmt.Errorf("not implemented")
}

func (p pClient) EvictV1beta1(ctx context.Context, eviction *policyv1beta1.Eviction) error {
	return fmt.Errorf("not implemented")
}

func (p pClient) GetLogs(name string, opts *v1.PodLogOptions) *restclient.Request {
	return nil
}

func (p pClient) ProxyGet(scheme, name, port, path string, params map[string]string) restclient.ResponseWrapper {
	return nil
}

func NewPClient() v1core.PodInterface {
	return pClient{}
}

func NewPod(ready bool, hostIP string) v1.Pod {
	podReadyStatus := v1.ConditionTrue
	if !ready {
		podReadyStatus = v1.ConditionFalse
	}

	return v1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: PodEnvVariableValue},
		Status: v1.PodStatus{
			HostIP: hostIP,
			Conditions: []v1.PodCondition{
				{
					Type:   v1.PodReady,
					Status: podReadyStatus,
				},
			},
			ContainerStatuses: []v1.ContainerStatus{
				{
					Name:  MockContainerName,
					Ready: ready,
				},
			},
		},
	}
}
