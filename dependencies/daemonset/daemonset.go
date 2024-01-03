package daemonset

import (
	"context"
	"fmt"
	"os"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	entry "opendev.org/airship/kubernetes-entrypoint/entrypoint"
	"opendev.org/airship/kubernetes-entrypoint/logger"
	"opendev.org/airship/kubernetes-entrypoint/util/env"
)

const (
	PodNameEnvVar            = "POD_NAME"
	PodNameNotSetErrorFormat = "env POD_NAME not set, daemonset dependency %s in namespace %s will be ignored"
)

type Daemonset struct {
	name      string
	namespace string
	podName   string
}

func init() {
	daemonsetEnv := fmt.Sprintf("%sDAEMONSET", entry.DependencyPrefix)
	daemonsetsDeps := env.SplitEnvToDeps(daemonsetEnv)
	for _, dep := range daemonsetsDeps {
		daemonset, err := NewDaemonset(dep.Name, dep.Namespace)
		if err != nil {
			logger.Error.Printf("Cannot initialize daemonset: %v", err)
			continue
		}
		entry.Register(daemonset)
	}
}

func NewDaemonset(name string, namespace string) (*Daemonset, error) {
	if os.Getenv(PodNameEnvVar) == "" {
		return nil, fmt.Errorf(PodNameNotSetErrorFormat, name, namespace)
	}
	return &Daemonset{
		name:      name,
		namespace: namespace,
		podName:   os.Getenv(PodNameEnvVar),
	}, nil
}

func (d Daemonset) IsResolved(ctx context.Context, entrypoint entry.EntrypointInterface) (bool, error) {
	var myPodName string
	daemonset, err := entrypoint.Client().DaemonSets(d.namespace).Get(ctx, d.name, metav1.GetOptions{})

	if err != nil {
		return false, err
	}

	label := metav1.FormatLabelSelector(daemonset.Spec.Selector)
	opts := metav1.ListOptions{LabelSelector: label}

	daemonsetPods, err := entrypoint.Client().Pods(d.namespace).List(ctx, opts)
	if err != nil {
		return false, err
	}

	myPod, err := entrypoint.Client().Pods(env.GetBaseNamespace()).Get(ctx, d.podName, metav1.GetOptions{})
	if err != nil {
		return false, fmt.Errorf("getting POD: %v failed : %v", myPodName, err)
	}

	myHost := myPod.Status.HostIP

	for _, pod := range daemonsetPods.Items {
		pod := pod // pinning
		if !isPodOnHost(&pod, myHost) {
			continue
		}
		if isPodReady(pod) {
			return true, nil
		}
		return false, fmt.Errorf("pod %v of daemonset %s is not ready", pod.Name, d)

	}
	return true, nil
}

func isPodOnHost(pod *v1.Pod, hostIP string) bool {
	return pod.Status.HostIP == hostIP
}

func isPodReady(pod v1.Pod) bool {
	for _, condition := range pod.Status.Conditions {
		if condition.Type == v1.PodReady && condition.Status == "True" {
			return true
		}
	}
	return false
}

func (d Daemonset) String() string {
	return fmt.Sprintf("Daemonset %s in namespace %s", d.name, d.namespace)
}
