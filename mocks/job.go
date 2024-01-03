package mocks

import (
	"context"
	"fmt"

	v1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	batchv1applyconfigurations "k8s.io/client-go/applyconfigurations/batch/v1"
	v1batch "k8s.io/client-go/kubernetes/typed/batch/v1"
)

const (
	SucceedingJobName  = "succeed"
	FailingJobName     = "fail"
	SucceedingJobLabel = "succeed"
	FailingJobLabel    = "fail"
)

type jClient struct {
}

func (j jClient) Create(ctx context.Context, job *v1.Job, opts metav1.CreateOptions) (*v1.Job, error) {
	return nil, fmt.Errorf("not implemented")
}

func (j jClient) Update(ctx context.Context, job *v1.Job, opts metav1.UpdateOptions) (*v1.Job, error) {
	return nil, fmt.Errorf("not implemented")
}

func (j jClient) UpdateStatus(ctx context.Context, job *v1.Job, opts metav1.UpdateOptions) (*v1.Job, error) {
	return nil, fmt.Errorf("not implemented")
}

func (j jClient) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return fmt.Errorf("not implemented")
}

func (j jClient) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	return fmt.Errorf("not implemented")
}

func (j jClient) Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.Job, error) {
	if name == SucceedingJobName {
		return &v1.Job{
			Status: v1.JobStatus{Succeeded: 1},
		}, nil
	}
	if name == FailingJobName {
		return &v1.Job{
			Status: v1.JobStatus{Succeeded: 0},
		}, nil
	}
	return nil, fmt.Errorf("mock job didnt work")
}

func (j jClient) List(ctx context.Context, opts metav1.ListOptions) (*v1.JobList, error) {
	var jobs []v1.Job

	switch opts.LabelSelector {
	case fmt.Sprintf("name=%s", SucceedingJobLabel):
		jobs = []v1.Job{NewJob(1)}
	case fmt.Sprintf("name=%s", FailingJobLabel):
		jobs = []v1.Job{NewJob(1), NewJob(0)}
	default:
		return nil, fmt.Errorf("mock job didnt work")
	}

	return &v1.JobList{
		Items: jobs,
	}, nil
}

func (j jClient) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return nil, fmt.Errorf("not implemented")
}

func (j jClient) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.Job, err error) {
	return nil, fmt.Errorf("not implemented")
}

func (j jClient) Apply(ctx context.Context, job *batchv1applyconfigurations.JobApplyConfiguration, opts metav1.ApplyOptions) (result *v1.Job, err error) {
	return nil, fmt.Errorf("not implemented")
}

func (j jClient) ApplyStatus(ctx context.Context, job *batchv1applyconfigurations.JobApplyConfiguration, opts metav1.ApplyOptions) (result *v1.Job, err error) {
	return nil, fmt.Errorf("not implemented")
}

func NewJClient() v1batch.JobInterface {
	return jClient{}
}

func NewJob(succeeded int32) v1.Job {
	return v1.Job{
		Status: v1.JobStatus{Succeeded: succeeded},
	}
}
