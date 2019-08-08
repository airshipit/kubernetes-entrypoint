package env

import (
	"os"
	"reflect"
	"testing"
)

const (
	name1            = "foo"
	name2            = "bar"
	defaultNamespace = "default"
	altNamespace1    = "fooNS"
	altNamespace2    = "barNS"
)

func TestSplitEnvToListWithComma(t *testing.T) {
	defer os.Unsetenv("TEST_LIST")
	os.Setenv("TEST_LIST", name1+","+name2)
	list := SplitEnvToDeps("TEST_LIST")
	if len(list) != 2 {
		t.Errorf("Expected len to be 2 not %d", len(list))
	}
	if list[0].Name != name1 {
		t.Errorf("Expected: %s got %s", name1, list[0])
	}
	if list[1].Name != name2 {
		t.Errorf("Expected: %s got %s", name2, list[1])
	}
}

func TestSplitEnvToList(t *testing.T) {
	defer os.Unsetenv("TEST_LIST")
	os.Setenv("TEST_LIST", name1)
	list := SplitEnvToDeps("TEST_LIST")
	if len(list) != 1 {
		t.Errorf("Expected len to be 1 not %d", len(list))
	}
	if list[0].Name != name1 {
		t.Errorf("Expected: %s got %s", name1, list[0])
	}
	if list[0].Namespace != defaultNamespace {
		t.Errorf("Expected: %s got %s", defaultNamespace, list[0].Namespace)
	}
}

func TestSplitEnvToListWithColon(t *testing.T) {
	defer os.Unsetenv("TEST_LIST")
	os.Setenv("TEST_LIST", altNamespace1+":"+name1)
	list := SplitEnvToDeps("TEST_LIST")
	if len(list) != 1 {
		t.Errorf("Expected len to be 1 not %d", len(list))
	}
	if list[0].Name != name1 {
		t.Errorf("Expected: %s got %s", name1, list[0].Name)
	}
	if list[0].Namespace != altNamespace1 {
		t.Errorf("Expected: %s got %s", altNamespace1, list[0].Namespace)
	}
}

func TestSplitEnvToListWithTooManyColons(t *testing.T) {
	// TODO(howell): This should probably expect an error
	defer os.Unsetenv("TEST_LIST")
	os.Setenv("TEST_LIST", "too:many:colons")
	list := SplitEnvToDeps("TEST_LIST")
	if len(list) != 0 {
		t.Errorf("Expected list to be empty")
	}
}

func TestSplitEnvToListWithColonsAndCommas(t *testing.T) {
	defer os.Unsetenv("TEST_LIST")
	os.Setenv("TEST_LIST", altNamespace1+":"+name1+","+altNamespace2+":"+name2)
	list := SplitEnvToDeps("TEST_LIST")
	if len(list) != 2 {
		t.Errorf("Expected len to be 2 not %d", len(list))
	}
	if list[0].Name != name1 {
		t.Errorf("Expected: %s got %s", name1, list[0].Name)
	}
	if list[0].Namespace != altNamespace1 {
		t.Errorf("Expected: %s got %s", altNamespace1, list[0].Namespace)
	}
	if list[1].Name != name2 {
		t.Errorf("Expected: %s got %s", name2, list[0].Name)
	}
	if list[1].Namespace != altNamespace2 {
		t.Errorf("Expected: %s got %s", altNamespace2, list[0].Namespace)
	}
}

func TestSplitEnvToListWithMissingNamespace(t *testing.T) {
	defer os.Unsetenv("TEST_LIST")
	os.Setenv("TEST_LIST", ":name")
	list := SplitEnvToDeps("TEST_LIST")
	if len(list) != 0 {
		t.Errorf("Invalid format, missing namespace in pod")
	}
}

func TestSplitEmptyEnvWithColon(t *testing.T) {
	defer os.Unsetenv("TEST_LIST")
	os.Setenv("TEST_LIST", "")
	list := SplitEnvToDeps("TEST_LIST")
	if len(list) != 0 {
		t.Errorf("Expected list to be empty")
	}
}

func TestSplitPodEnvToDepsSuccess(t *testing.T) {
	testListJSONVal := `[
  {
    "namespace": "` + name1 + `",
    "labels": {
      "k1": "v1",
      "k2": "v2"
    },
    "requireSameNode": true
  },
  {
    "labels": {
      "k1": "v1",
      "k2": "v2"
     }
  }
]`
	os.Setenv("NAMESPACE", "TEST_NAMESPACE")
	os.Setenv("TEST_LIST_JSON", testListJSONVal)
	actual := SplitPodEnvToDeps("TEST_LIST_JSON")
	expected := []PodDependency{
		{
			Namespace: name1,
			Labels: map[string]string{
				"k1": "v1",
				"k2": "v2",
			},
			RequireSameNode: true,
		},
		{
			Namespace: "TEST_NAMESPACE",
			Labels: map[string]string{
				"k1": "v1",
				"k2": "v2",
			},
			RequireSameNode: false,
		},
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected: %v Got: %v", expected, actual)
	}
}

func TestSplitPodEnvToDepsUnset(t *testing.T) {
	defer os.Unsetenv("TEST_LIST")
	os.Setenv("TEST_LIST", "")
	actual := SplitPodEnvToDeps("TEST_LIST")
	if len(actual) != 0 {
		t.Errorf("Expected: no dependencies Got: %v", actual)
	}
}

func TestSplitPodEnvToDepsIgnoreInvalid(t *testing.T) {
	defer os.Unsetenv("TEST_LIST")
	os.Setenv("TEST_LIST", `[{"invalid": json}`)
	actual := SplitPodEnvToDeps("TEST_LIST")
	if len(actual) != 0 {
		t.Errorf("Expected: ignore invalid dependencies Got: %v", actual)
	}
}

func TestSplitJobEnvToDepsJsonSuccess(t *testing.T) {
	testListJSONVal := `[
  {
    "namespace": "` + altNamespace1 + `",
    "labels": {
      "k1": "v1",
      "k2": "v2"
    }
  },
  {
    "name": "` + name1 + `"
  }
]`
	defer os.Unsetenv("NAMESPACE")
	os.Setenv("NAMESPACE", "TEST_NAMESPACE")
	defer os.Unsetenv("TEST_LIST_JSON")
	os.Setenv("TEST_LIST_JSON", testListJSONVal)
	actual := SplitJobEnvToDeps("TEST_LIST", "TEST_LIST_JSON")
	expected := []JobDependency{
		{
			Namespace: altNamespace1,
			Labels: map[string]string{
				"k1": "v1",
				"k2": "v2",
			},
		},
		{
			Name:      name1,
			Namespace: "TEST_NAMESPACE",
		},
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected: %v Got: %v", expected, actual)
	}
}

func TestSplitJobEnvToDepsPlainSuccess(t *testing.T) {
	defer os.Unsetenv("NAMESPACE")
	os.Setenv("NAMESPACE", "TEST_NAMESPACE")
	defer os.Unsetenv("TEST_LIST")
	os.Setenv("TEST_LIST", "plain")
	actual := SplitJobEnvToDeps("TEST_LIST", "TEST_LIST_JSON")
	expected := []JobDependency{
		{
			Name:      "plain",
			Namespace: "TEST_NAMESPACE",
		},
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected: %v Got: %v", expected, actual)
	}
}

func TestSplitJobEnvToDepsJsonPrecedence(t *testing.T) {
	defer os.Unsetenv("NAMESPACE")
	os.Setenv("NAMESPACE", "TEST_NAMESPACE")
	defer os.Unsetenv("TEST_LIST_JSON")
	os.Setenv("TEST_LIST_JSON", `[{"name": "json"}]`)
	defer os.Unsetenv("TEST_LIST")
	os.Setenv("TEST_LIST", "plain")
	actual := SplitJobEnvToDeps("TEST_LIST", "TEST_LIST_JSON")
	expected := []JobDependency{
		{
			Name:      "json",
			Namespace: "TEST_NAMESPACE",
		},
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected: %v Got: %v", expected, actual)
	}
}

func TestSplitJobEnvToDepsUnset(t *testing.T) {
	actual := SplitJobEnvToDeps("TEST_LIST", "TEST_LIST_JSON")
	if len(actual) != 0 {
		t.Errorf("Expected: no dependencies Got: %v", actual)
	}
}

func TestSplitJobEnvToDepsIgnoreInvalid(t *testing.T) {
	defer os.Unsetenv("TEST_LIST_JSON")
	os.Setenv("TEST_LIST_JSON", `[{"invalid": json}`)
	actual := SplitJobEnvToDeps("TEST_LIST", "TEST_LIST_JSON")
	if len(actual) != 0 {
		t.Errorf("Expected: ignore invalid dependencies Got: %v", actual)
	}
}

func TestSplitCommandUnset(t *testing.T) {
	defer os.Unsetenv("COMMAND")
	list := SplitCommand()
	if len(list) > 0 {
		t.Errorf("Expected len to be 0, got %d", len(list))
	}
}

func TestSplitCommandEmpty(t *testing.T) {
	defer os.Unsetenv("COMMAND")
	os.Setenv("COMMAND", "")
	list := SplitCommand()
	if len(list) > 0 {
		t.Errorf("Expected len to be 0, got %v", len(list))
	}
}

func TestSplitCommand(t *testing.T) {
	defer os.Unsetenv("COMMAND")
	os.Setenv("COMMAND", "echo test")
	list := SplitCommand()
	if len(list) != 2 {
		t.Errorf("Expected two elements, got %v", len(list))
	}
	if list[0] != "echo" {
		t.Errorf("Expected echo, got %s", list[0])
	}
	if list[1] != "test" {
		t.Errorf("Expected test, got %s", list[1])
	}
}

func TestGetBaseNamespaceEmpty(t *testing.T) {
	defer os.Unsetenv("NAMESPACE")
	os.Setenv("NAMESPACE", "")
	getBaseNamespace := GetBaseNamespace()
	if getBaseNamespace != defaultNamespace {
		t.Errorf("Expected namespace to be %s, got %s", defaultNamespace, getBaseNamespace)
	}
}

func TestGetBaseNamespace(t *testing.T) {
	defer os.Unsetenv("NAMESPACE")
	os.Setenv("NAMESPACE", "foo")
	getBaseNamespace := GetBaseNamespace()
	if getBaseNamespace != "foo" {
		t.Errorf("Expected namespace to be foo, got %v", getBaseNamespace)
	}
}
