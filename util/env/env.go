package env

import (
	"encoding/json"
	"os"
	"strings"

	"opendev.org/airship/kubernetes-entrypoint/logger"
)

const (
	Separator = ":"
)

type Dependency struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type PodDependency struct {
	Labels          map[string]string `json:"labels"`
	Namespace       string            `json:"namespace"`
	RequireSameNode bool              `json:"requireSameNode"`
}

type JobDependency struct {
	Name      string            `json:"name"`
	Labels    map[string]string `json:"labels"`
	Namespace string            `json:"namespace"`
}

type CustomResourceDependency struct {
	APIVersion string              `json:"apiVersion"`
	Name       string              `json:"name"`
	Namespace  string              `json:"namespace"`
	Kind       string              `json:"kind"`
	Fields     []map[string]string `json:"fields"`
}

func SplitCommand() []string {
	command := os.Getenv("COMMAND")
	if command == "" {
		return []string{}
	}
	commandList := strings.Split(command, " ")
	return commandList
}

// SplitEnvToDeps returns list of namespaces and names pairs
func SplitEnvToDeps(env string) (envList []Dependency) {
	separator := ","

	e := os.Getenv(env)
	if e == "" {
		return envList
	}

	envVars := strings.Split(e, separator)
	namespace := GetBaseNamespace()
	var dep Dependency
	for _, envVar := range envVars {
		if strings.Contains(envVar, Separator) {
			nameAfterSplit := strings.Split(envVar, Separator)
			if len(nameAfterSplit) != 2 {
				logger.Warning.Printf("Invalid format got %s, expected namespace:name", envVar)
				continue
			}
			if nameAfterSplit[0] == "" {
				logger.Warning.Printf("Invalid format, missing namespace %s", envVar)
				continue
			}
			dep = Dependency{Name: nameAfterSplit[1], Namespace: nameAfterSplit[0]}
		} else {
			dep = Dependency{Name: envVar, Namespace: namespace}
		}
		envList = append(envList, dep)
	}

	return envList
}

// SplitPodEnvToDeps returns list of PodDependency
func SplitPodEnvToDeps(env string) []PodDependency {
	deps := []PodDependency{}

	namespace := GetBaseNamespace()

	e := os.Getenv(env)
	if e == "" {
		return deps
	}

	err := json.Unmarshal([]byte(e), &deps)
	if err != nil {
		logger.Warning.Printf("Invalid format: %v", e)
		return []PodDependency{}
	}

	for i, dep := range deps {
		if dep.Namespace == "" {
			dep.Namespace = namespace
		}
		deps[i] = dep
	}

	return deps
}

// SplitJobEnvToDeps returns list of JobDependency
func SplitJobEnvToDeps(env string, jsonEnv string) []JobDependency {
	deps := []JobDependency{}

	namespace := GetBaseNamespace()

	envVal := os.Getenv(env)
	jsonEnvVal := os.Getenv(jsonEnv)
	if jsonEnvVal != "" {
		if envVal != "" {
			logger.Warning.Printf("Ignoring %s since %s was specified", env, jsonEnv)
		}
		err := json.Unmarshal([]byte(jsonEnvVal), &deps)
		if err != nil {
			logger.Warning.Printf("Invalid format: %s", jsonEnvVal)
			return []JobDependency{}
		}

		valid := []JobDependency{}
		for _, dep := range deps {
			if dep.Namespace == "" {
				dep.Namespace = namespace
			}

			valid = append(valid, dep)
		}

		return valid
	}

	if envVal != "" {
		plainDeps := SplitEnvToDeps(env)

		deps = []JobDependency{}
		for _, dep := range plainDeps {
			deps = append(deps, JobDependency{Name: dep.Name, Namespace: dep.Namespace})
		}

		return deps
	}

	return deps
}

// GetBaseNamespace returns default namespace when user set empty one
func GetBaseNamespace() string {
	namespace, isSet := os.LookupEnv("NAMESPACE")
	if !isSet || namespace == "" {
		namespace = "default"
	}
	return namespace
}
