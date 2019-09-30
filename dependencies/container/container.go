package container

import (
	"fmt"
	"os"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	entry "opendev.org/airship/kubernetes-entrypoint/entrypoint"
	"opendev.org/airship/kubernetes-entrypoint/logger"
	"opendev.org/airship/kubernetes-entrypoint/util"
	"opendev.org/airship/kubernetes-entrypoint/util/env"
)

const (
	PodNameNotSetError    = "Environment variable POD_NAME not set"
	NamespaceNotSupported = "Container doesn't accept namespace"
)

type Container struct {
	name string
}

func init() {
	containerEnv := fmt.Sprintf("%sCONTAINER", entry.DependencyPrefix)
	if util.ContainsSeparator(containerEnv, "Container") {
		logger.Error.Printf(NamespaceNotSupported)
		os.Exit(1)
	}
	if containerDeps := env.SplitEnvToDeps(containerEnv); containerDeps != nil {

		if len(containerDeps) > 0 {
			for _, dep := range containerDeps {
				entry.Register(NewContainer(dep.Name))
			}
		}
	}
}

func NewContainer(name string) Container {
	return Container{name: name}

}

func (c Container) IsResolved(entrypoint entry.EntrypointInterface) (bool, error) {
	myPodName := os.Getenv("POD_NAME")
	if myPodName == "" {
		return false, fmt.Errorf(PodNameNotSetError)
	}
	pod, err := entrypoint.Client().Pods(env.GetBaseNamespace()).Get(myPodName, metav1.GetOptions{})
	if err != nil {
		return false, err
	}

	if strings.Contains(c.name, env.Separator) {
		return false, fmt.Errorf("Specifying namespace is not permitted")
	}
	containers := pod.Status.ContainerStatuses
	for _, container := range containers {
		if container.Name == c.name && container.Ready {
			return true, nil
		}
	}
	return false, nil
}

func (c Container) String() string {
	return fmt.Sprintf("Container %s", c.name)
}
