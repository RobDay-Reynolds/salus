package main_test

import (
	"encoding/json"
	"fmt"
	"github.com/cznic/fileutil"
	"github.com/golang/go/src/pkg/io/ioutil"
	"github.com/monkeyherder/moirai/checks/network"
	. "github.com/monkeyherder/moirai/main"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"os"
	"os/exec"
	"syscall"
	"time"
)

var _ = Describe("Checksd", func() {
	AfterEach(func() {
		gexec.KillAndWait(10)
	})

	Context("Given valid config", func() {
		var command *exec.Cmd

		BeforeEach(func() {
			config := ChecksdConfig{
				ChecksPollTime: 1 * time.Second,
				IcmpChecks: []network.IcmpCheck{
					{
						Address: "www.google.com",
						Timeout: 5 * time.Second,
					},
				},
			}

			configJson, err := json.Marshal(config)
			Expect(err).ToNot(HaveOccurred())

			tempConfigFile, err := fileutil.TempFile(os.TempDir(), "checksd", "yml")
			Expect(err).ToNot(HaveOccurred())

			err = ioutil.WriteFile(tempConfigFile.Name(), configJson, os.ModePerm)
			Expect(err).ToNot(HaveOccurred())

			command = exec.Command(pathToChecksd, fmt.Sprintf("--config=%s", tempConfigFile.Name()))
		})

		It("Should run health check", func() {
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).ShouldNot(HaveOccurred())
			Eventually(session, 10).Should(gbytes.Say("After Check ran"))
		})

		Context("Sending a SIGTERM process signal", func() {
			It("should handle signal and exit gracefully", func() {
				session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).ShouldNot(HaveOccurred())
				Eventually(session, 10).Should(gbytes.Say("After Check ran"))

				session.Signal(syscall.SIGTERM)
				Eventually(session.ExitCode, 3).Should(Equal(0))
				Eventually(session, 10).Should(gbytes.Say("sig received: terminated"))
			})
		})
	})

	Context("Given invalid config", func() {
		var command *exec.Cmd

		BeforeEach(func() {
			tempConfigFile, err := fileutil.TempFile(os.TempDir(), "checksd", "yml")
			Expect(err).ToNot(HaveOccurred())

			err = ioutil.WriteFile(tempConfigFile.Name(), []byte(`not json`), os.ModePerm)
			Expect(err).ToNot(HaveOccurred())

			command = exec.Command(pathToChecksd, fmt.Sprintf("--config=%s", tempConfigFile.Name()))
		})

		It("Should error", func() {
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).ShouldNot(HaveOccurred())
			Eventually(session.Err, 10).Should(gbytes.Say("unable to configure checksd with config file"))
		})
	})
})