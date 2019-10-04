package mocks

import (
	"opendev.org/airship/kubernetes-entrypoint/client"
)

type MockEntrypoint struct {
	MockClient *Client
	namespace  string
}

func (m MockEntrypoint) Resolve() {}

func (m MockEntrypoint) Client() client.ClientInterface {
	return m.MockClient
}

func (m MockEntrypoint) GetNamespace() string {
	return m.namespace
}

func NewEntrypointInNamespace(namespace string) *MockEntrypoint {
	return &MockEntrypoint{
		MockClient: NewClient(),
		namespace:  namespace,
	}
}

func NewEntrypoint() *MockEntrypoint {
	return NewEntrypointInNamespace("test")
}
