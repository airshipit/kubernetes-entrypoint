package daemonset_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDaemonset(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Daemonset Suite")
}
