package socket

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"opendev.org/airship/kubernetes-entrypoint/entrypoint"
	"opendev.org/airship/kubernetes-entrypoint/mocks"
)

const (
	tempPathSuffix = "k8s-entrypoint"

	existingSocket    = "existing-socket"
	nonExistingSocket = "nonexisting-socket"
)

var (
	testDir string

	existingSocketPath    string
	nonExistingSocketPath string
)

var testEntrypoint entrypoint.EntrypointInterface

var _ = Describe("Socket", func() {

	// NOTE: It is impossible for a user to create a file that he does not
	// have access to, and thus it is impossible to write an isolated unit
	// test that checks for permission errors. That test is omitted from
	// this suite

	BeforeEach(func() {
		testEntrypoint = mocks.NewEntrypoint()

		var err error
		testDir, err = ioutil.TempDir("", tempPathSuffix)
		Expect(err).NotTo(HaveOccurred())

		existingSocketPath = filepath.Join(testDir, existingSocket)
		nonExistingSocketPath = filepath.Join(testDir, nonExistingSocket)

		_, err = os.Create(existingSocketPath)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		err := os.RemoveAll(testDir)
		Expect(err).NotTo(HaveOccurred())
	})

	It("checks the name of a newly created socket", func() {
		socket := NewSocket(existingSocketPath)

		Expect(socket.name).To(Equal(existingSocketPath))
	})

	It("resolves an existing socket socket", func() {
		socket := NewSocket(existingSocketPath)

		isResolved, err := socket.IsResolved(testEntrypoint)

		Expect(isResolved).To(Equal(true))
		Expect(err).NotTo(HaveOccurred())
	})

	It("fails on trying to resolve a nonexisting socket", func() {
		socket := NewSocket(nonExistingSocketPath)

		isResolved, err := socket.IsResolved(testEntrypoint)

		Expect(isResolved).To(Equal(false))
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal(fmt.Sprintf(NonExistingErrorFormat, socket)))
	})
})
