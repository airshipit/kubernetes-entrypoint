package config

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"opendev.org/airship/kubernetes-entrypoint/entrypoint"
	"opendev.org/airship/kubernetes-entrypoint/mocks"
)

const (
	testDir        = "/tmp"
	interfaceName  = "INTERFACE_NAME"
	testConfigName = "KUBERNETES_ENTRYPOINT_TEST_CONFIG"

	testConfigContentsFormat = "TEST_CONFIG %s\n"

	templatePrefix = "/tmp/templates"
)

var testEntrypoint entrypoint.EntrypointInterface
var testConfigContents string
var testConfigPath string
var testTemplatePath string

// var testClient cli.ClientInterface

func init() {
	testConfigContents = fmt.Sprintf(testConfigContentsFormat, "{{ .HOSTNAME }}")

	testTemplatePath = fmt.Sprintf("%s/%s/%s", templatePrefix, testConfigName, testConfigName)
	testConfigPath = fmt.Sprintf("%s/%s", testDir, testConfigName)
}

func setupOsEnvironment() (err error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return err
	}

	ifaceName := ifaces[0].Name
	return os.Setenv(interfaceName, ifaceName)
}

func teardownOsEnvironment() (err error) {
	return os.Unsetenv(interfaceName)
}

func setupConfigTemplate(templatePath string) error {
	configContent := []byte(testConfigContents)
	if err := os.MkdirAll(filepath.Dir(templatePath), 0755); err != nil {
		return err
	}

	if err := ioutil.WriteFile(templatePath, configContent, 0644); err != nil {
		return err
	}

	return nil
}

func teardownConfigTemplate() (err error) {
	if err := os.RemoveAll(templatePrefix); err != nil {
		return err
	}

	if err := os.RemoveAll(testConfigPath); err != nil {
		return err
	}

	return
}

var _ = Describe("Config", func() {

	BeforeEach(func() {
		err := setupOsEnvironment()
		Expect(err).NotTo(HaveOccurred())

		err = setupConfigTemplate(testTemplatePath)
		Expect(err).NotTo(HaveOccurred())

		testEntrypoint = mocks.NewEntrypoint()
	})

	AfterEach(func() {
		err := teardownOsEnvironment()
		Expect(err).NotTo(HaveOccurred())

		err = teardownConfigTemplate()
		Expect(err).NotTo(HaveOccurred())
	})

	It("creates new config from file", func() {
		config, err := NewConfig(testConfigPath, templatePrefix)

		Expect(config).NotTo(Equal(nil))
		Expect(err).NotTo(HaveOccurred())
	})

	It("checks the name of a newly created config file", func() {
		config, err := NewConfig(testConfigPath, templatePrefix)
		Expect(config.name).To(Equal(testConfigPath))

		Expect(config).NotTo(Equal(nil))
		Expect(err).NotTo(HaveOccurred())
	})

	It("checks the format of a newly created config file", func() {
		config, _ := NewConfig(testConfigPath, templatePrefix)
		_, err := config.IsResolved(testEntrypoint)
		Expect(err).NotTo(HaveOccurred())

		result, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", testDir, testConfigName))
		Expect(err).NotTo(HaveOccurred())

		hostname, err := os.Hostname()
		Expect(err).NotTo(HaveOccurred())

		expectedFile := fmt.Sprintf(testConfigContentsFormat, hostname)

		readConfig := string(result)
		Expect(readConfig).To(BeEquivalentTo(expectedFile))
	})

	It("checks resolution of a config", func() {
		config, _ := NewConfig(testConfigPath, templatePrefix)

		isResolved, err := config.IsResolved(testEntrypoint)

		Expect(isResolved).To(Equal(true))
		Expect(err).NotTo(HaveOccurred())
	})
})
