package customresource

import (
	"context"
	"errors"
	"os"
	"reflect"
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"opendev.org/airship/kubernetes-entrypoint/mocks"
)

const listOfTwoCRResolvers = `
[
  {
    "apiVersion": "api1",
    "kind": "kind1",
    "namespace": "foospace1",
    "name": "foo1",
    "fields": [
      {
        "key": "field1key1",
        "value": "field1val1"
      },
      {
        "key": "field1key2",
        "value": "field1val2"
      }
    ]
  },
  {
    "apiVersion": "api2",
    "kind": "kind2",
    "namespace": "foospace2",
    "name": "foo2",
    "fields": [
      {
        "key": "field2key1",
        "value": "field2val1"
      },
      {
        "key": "field2key2",
        "value": "field2val2"
      }
    ]
  }
]`

const listOfConditionResolvers = `
[
  {
    "apiVersion": "api1",
    "kind": "kind1",
    "namespace": "foospace1",
    "name": "foo1",
    "conditions": [
      {
        "type": "Ready",
        "status": "True"
      },
      {
        "type": "Degraded",
        "status": "False"
      }
    ]
  },
  {
    "apiVersion": "api2",
    "kind": "kind2",
    "namespace": "foospace2",
    "name": "foo2",
    "conditions": [
      {
        "type": "Ready"
      }
    ]
  }
]`

// conditioned builds a resource whose status carries the given conditions,
// each a type/status pair.
func conditioned(conditions ...[2]string) *unstructured.Unstructured {
	items := []interface{}{}
	for _, condition := range conditions {
		items = append(items, map[string]interface{}{
			"type":               condition[0],
			"status":             condition[1],
			"reason":             "Reconciled",
			"message":            "",
			"lastTransitionTime": "2026-08-08T00:00:00Z",
		})
	}
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "stable.example.com/v1",
			"kind":       "Foo",
			"name":       "my-foo",
			"namespace":  "default",
			"status": map[string]interface{}{
				"conditions": items,
			},
		},
	}
}

func TestFromEnv(t *testing.T) {
	tests := []struct {
		name      string
		useEnvVar bool
		envVar    string
		expected  []Resolver
		expectErr bool
	}{
		{
			name:      "EmptyVar",
			useEnvVar: false,
			envVar:    "",
			expected:  []Resolver{},
			expectErr: false,
		},
		{
			name:      "InvalidVar",
			useEnvVar: true,
			envVar:    `[}"invalid": "json"}]`,
			expected:  []Resolver{},
			expectErr: true,
		},
		{
			name:      "Successful",
			useEnvVar: true,
			envVar:    listOfTwoCRResolvers,
			expected: []Resolver{
				{
					APIVersion: "api1",
					Kind:       "kind1",
					Namespace:  "foospace1",
					Name:       "foo1",
					Fields: []Field{
						{
							Key:   "field1key1",
							Value: "field1val1",
						},
						{
							Key:   "field1key2",
							Value: "field1val2",
						},
					},
				},
				{
					APIVersion: "api2",
					Kind:       "kind2",
					Namespace:  "foospace2",
					Name:       "foo2",
					Fields: []Field{
						{
							Key:   "field2key1",
							Value: "field2val1",
						},
						{
							Key:   "field2key2",
							Value: "field2val2",
						},
					},
				},
			},
			expectErr: false,
		},
		{
			name:      "UnNamespaced",
			useEnvVar: true,
			envVar:    `[{"apiVersion":"api","kind":"kind","name":"foo","fields":[]}]`,
			expected: []Resolver{
				{
					APIVersion: "api",
					Kind:       "kind",
					Namespace:  "default",
					Name:       "foo",
					Fields:     []Field{},
				},
			},
			expectErr: false,
		},
		{
			name:      "Conditions",
			useEnvVar: true,
			envVar:    listOfConditionResolvers,
			expected: []Resolver{
				{
					APIVersion: "api1",
					Kind:       "kind1",
					Namespace:  "foospace1",
					Name:       "foo1",
					Conditions: []Condition{
						{
							Type:   "Ready",
							Status: "True",
						},
						{
							Type:   "Degraded",
							Status: "False",
						},
					},
				},
				{
					APIVersion: "api2",
					Kind:       "kind2",
					Namespace:  "foospace2",
					Name:       "foo2",
					Conditions: []Condition{
						{
							Type: "Ready",
						},
					},
				},
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.useEnvVar {
				os.Setenv("TEST_LIST_JSON", tt.envVar)
				defer os.Unsetenv("TEST_LIST_JSON")
			} else {
				os.Unsetenv("TEST_LIST_JSON")
			}

			actual, err := fromEnv("TEST_LIST_JSON")
			if err != nil {
				if !tt.expectErr {
					t.Fatalf("Unexpected error: %s", err.Error())
				}
			} else if tt.expectErr {
				t.Fatalf("Expected error, but received none")
			}

			if !reflect.DeepEqual(tt.expected, actual) {
				t.Errorf("Expected %+v, got %+v", tt.expected, actual)
			}
		})
	}
}

func TestIsResolved(t *testing.T) {
	tests := []struct {
		name           string
		customResource *unstructured.Unstructured
		resolver       Resolver
		expected       bool
		expectErr      bool
		clientErr      error
	}{
		{
			name: "Successful",
			customResource: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "stable.example.com/v1",
					"kind":       "Foo",
					"name":       "my-foo",
					"namespace":  "default",
					"simpleKey":  "ready",
					"complex": map[string]interface{}{
						"key": map[string]interface{}{
							"with": map[string]interface{}{
								"layers": "ready",
							},
						},
					},
				},
			},
			resolver: Resolver{
				APIVersion: "stable.exmaple.com/v1",
				Kind:       "Foo",
				Name:       "my-foo",
				Fields: []Field{
					{
						Key:   "simpleKey",
						Value: "ready",
					},
					{
						Key:   "complex.key.with.layers",
						Value: "ready",
					},
				},
			},
			expected:  true,
			expectErr: false,
			clientErr: nil,
		},
		{
			name: "Unresolved",
			customResource: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "stable.example.com/v1",
					"kind":       "Foo",
					"name":       "my-foo",
					"namespace":  "default",
					"simpleKey":  "NOT-ready",
				},
			},
			resolver: Resolver{
				APIVersion: "stable.exmaple.com/v1",
				Kind:       "Foo",
				Name:       "my-foo",
				Fields: []Field{
					{
						Key:   "simpleKey",
						Value: "ready",
					},
				},
			},
			expected:  false,
			expectErr: true,
			clientErr: nil,
		},
		{
			name:           "ClientError",
			customResource: &unstructured.Unstructured{},
			resolver:       Resolver{},
			expected:       false,
			expectErr:      true,
			clientErr:      errors.New("Generic error"),
		},
		{
			name: "BadCustomResource",
			customResource: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "stable.example.com/v1",
					"kind":       "Foo",
					"name":       "my-foo",
					"namespace":  "default",
					"notAMap":    5,
				},
			},
			resolver: Resolver{
				APIVersion: "stable.exmaple.com/v1",
				Kind:       "Foo",
				Name:       "my-foo",
				Fields: []Field{
					{
						Key:   "notAMap.nothingHere",
						Value: "ready",
					},
				},
			},
			expected:  false,
			expectErr: true,
			clientErr: nil,
		},
		{
			name: "MissingKey",
			customResource: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "stable.example.com/v1",
					"kind":       "Foo",
					"name":       "my-foo",
					"namespace":  "default",
				},
			},
			resolver: Resolver{
				APIVersion: "stable.exmaple.com/v1",
				Kind:       "Foo",
				Name:       "my-foo",
				Fields: []Field{
					{
						Key:   "missingKey",
						Value: "ready",
					},
				},
			},
			expected:  false,
			expectErr: true,
			clientErr: nil,
		},
		{
			name:           "ConditionReady",
			customResource: conditioned([2]string{"Ready", "True"}),
			resolver: Resolver{
				Kind:       "Foo",
				Name:       "my-foo",
				Conditions: []Condition{{Type: "Ready", Status: "True"}},
			},
			expected:  true,
			expectErr: false,
		},
		{
			name:           "ConditionStatusDefaultsToTrue",
			customResource: conditioned([2]string{"Ready", "True"}),
			resolver: Resolver{
				Kind:       "Foo",
				Name:       "my-foo",
				Conditions: []Condition{{Type: "Ready"}},
			},
			expected:  true,
			expectErr: false,
		},
		{
			name:           "ConditionNotReadyYet",
			customResource: conditioned([2]string{"Ready", "False"}),
			resolver: Resolver{
				Kind:       "Foo",
				Name:       "my-foo",
				Conditions: []Condition{{Type: "Ready"}},
			},
			expected:  false,
			expectErr: true,
		},
		{
			name:           "ConditionUnknown",
			customResource: conditioned([2]string{"Ready", "Unknown"}),
			resolver: Resolver{
				Kind:       "Foo",
				Name:       "my-foo",
				Conditions: []Condition{{Type: "Ready"}},
			},
			expected:  false,
			expectErr: true,
		},
		{
			name:           "ConditionSelectedByType",
			customResource: conditioned([2]string{"Degraded", "False"}, [2]string{"Ready", "True"}),
			resolver: Resolver{
				Kind: "Foo",
				Name: "my-foo",
				Conditions: []Condition{
					{Type: "Ready"},
					{Type: "Degraded", Status: "False"},
				},
			},
			expected:  true,
			expectErr: false,
		},
		{
			name:           "ConditionTypeAbsent",
			customResource: conditioned([2]string{"Degraded", "False"}),
			resolver: Resolver{
				Kind:       "Foo",
				Name:       "my-foo",
				Conditions: []Condition{{Type: "Ready"}},
			},
			expected:  false,
			expectErr: true,
		},
		{
			name: "NoConditionsAtAll",
			customResource: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "stable.example.com/v1",
					"kind":       "Foo",
					"name":       "my-foo",
				},
			},
			resolver: Resolver{
				Kind:       "Foo",
				Name:       "my-foo",
				Conditions: []Condition{{Type: "Ready"}},
			},
			expected:  false,
			expectErr: true,
		},
		{
			name: "ConditionsNotAList",
			customResource: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "stable.example.com/v1",
					"kind":       "Foo",
					"name":       "my-foo",
					"status":     map[string]interface{}{"conditions": "not-a-list"},
				},
			},
			resolver: Resolver{
				Kind:       "Foo",
				Name:       "my-foo",
				Conditions: []Condition{{Type: "Ready"}},
			},
			expected:  false,
			expectErr: true,
		},
		{
			name: "ConditionStatusNotAString",
			customResource: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "stable.example.com/v1",
					"kind":       "Foo",
					"name":       "my-foo",
					"status": map[string]interface{}{
						"conditions": []interface{}{
							map[string]interface{}{"type": "Ready", "status": true},
						},
					},
				},
			},
			resolver: Resolver{
				Kind:       "Foo",
				Name:       "my-foo",
				Conditions: []Condition{{Type: "Ready"}},
			},
			expected:  false,
			expectErr: true,
		},
		{
			name:           "FieldsAndConditionsTogether",
			customResource: conditioned([2]string{"Ready", "True"}),
			resolver: Resolver{
				Kind:       "Foo",
				Name:       "my-foo",
				Fields:     []Field{{Key: "kind", Value: "Foo"}},
				Conditions: []Condition{{Type: "Ready"}},
			},
			expected:  true,
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ep := mocks.NewEntrypoint()
			ep.MockClient.FakeCustomResource = tt.customResource
			ep.MockClient.Err = tt.clientErr
			result, err := tt.resolver.IsResolved(context.TODO(), ep)
			if err != nil {
				if !tt.expectErr {
					t.Fatalf("Unexpected error: %s", err.Error())
				}
			} else if tt.expectErr {
				t.Fatalf("Expected error, but received none")
			}

			if result != tt.expected {
				t.Errorf("Expected success to be %v, but got %v", tt.expected, result)
			}
		})
	}
}
