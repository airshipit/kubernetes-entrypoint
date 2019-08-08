package entrypoint_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestEntrypoint(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Entrypoint Suite")
}
