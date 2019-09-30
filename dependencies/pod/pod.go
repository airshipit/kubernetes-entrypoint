package pod

import (
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
	PodNameNotSetErrorFormat = "Env POD_NAME not set. Pod dependency in namespace %s will be ignored!"
)

type Pod struct {
	namespace       string
	labels          map[string]string
	requireSameNode bool
	podName         string
}

func init() {
	podEnv := fmt.Sprintf("%sPOD%s", entry.DependencyPrefix, entry.JsonSuffix)
	podDeps := env.SplitPodEnvToDeps(podEnv)
	for _, dep := range podDeps {
		pod, err := NewPod(dep.Labels, dep.Namespace, dep.RequireSameNode)
		if err != nil {
			logger.Error.Printf("Cannot initialize pod: %v", err)
			continue
		}
		entry.Register(pod)
	}
}

func NewPod(labels map[string]string, namespace string, requireSameNode bool) (*Pod, error) {
	if os.Getenv(PodNameEnvVar) == "" {
		return nil, fmt.Errorf(PodNameNotSetErrorFormat, namespace)
	}
	return &Pod{
		labels:          labels,
		namespace:       namespace,
		requireSameNode: requireSameNode,
		podName:         os.Getenv(PodNameEnvVar),
	}, nil
}

func (p Pod) IsResolved(entrypoint entry.EntrypointInterface) (bool, error) {
	myPod, err := entrypoint.Client().Pods(env.GetBaseNamespace()).Get(p.podName, metav1.GetOptions{})
	if err != nil {
		return false, fmt.Errorf("Getting POD: %v failed : %v", p.podName, err)
	}
	myHost := myPod.Status.HostIP

	labelSelector := &metav1.LabelSelector{MatchLabels: p.labels}
	label := metav1.FormatLabelSelector(labelSelector)
	opts := metav1.ListOptions{LabelSelector: label}

	matchingPodList, err := entrypoint.Client().Pods(p.namespace).List(opts)
	if err != nil {
		return false, err
	}

	matchingPods := matchingPodList.Items
	if len(matchingPods) == 0 {
		return false, fmt.Errorf("Found no pods matching labels: %v", p.labels)
	}

	podCount := 0
	for _, pod := range matchingPods {
		podCount++
		pod := pod // pinning
		if p.requireSameNode && !isPodOnHost(&pod, myHost) {
			continue
		}
		if isPodReady(pod) {
			return true, nil
		}
	}
	onHostClause := ""
	if p.requireSameNode {
		onHostClause = " on host"
	}
	if podCount == 0 {
		return false, fmt.Errorf("Found no pods%v matching labels: %v", onHostClause, p.labels)
	} else {
		return false, fmt.Errorf("Found %v pods%v, but none ready, matching labels: %v", podCount, onHostClause, p.labels)
	}
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

func (p Pod) String() string {
	return fmt.Sprintf("Pod on same host with labels %v in namespace %s", p.labels, p.namespace)
}
