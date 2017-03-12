package main_test

import (
	"encoding/json"
	"fmt"
	"github.com/cznic/fileutil"
	"github.com/golang/go/src/pkg/path/filepath"
	"github.com/monkeyherder/moirai/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"syscall"
	"time"
)

var _ = Describe("Checksd", func() {

	act := func(config config.ChecksdConfig) *exec.Cmd {
		configJson, err := json.Marshal(config)
		Expect(err).ToNot(HaveOccurred())

		tempConfigFile, err := fileutil.TempFile(os.TempDir(), "checksd", "json")
		Expect(err).ToNot(HaveOccurred())

		err = ioutil.WriteFile(tempConfigFile.Name(), configJson, os.ModePerm)
		Expect(err).ToNot(HaveOccurred())
		return exec.Command(pathToChecksd, fmt.Sprintf("--config=%s", tempConfigFile.Name()))
	}

	AfterEach(func() {
		gexec.KillAndWait(10)
	})

	Context("Given valid config", func() {
		var checksdConfig config.ChecksdConfig
		var fileSummaryJsonPath string

		BeforeEach(func() {
			fileSummaryJsonPath = filepath.Join(os.TempDir(), "summary.json")

			checksdConfig = config.ChecksdConfig{
				CheckStatusFilePath: fileSummaryJsonPath,
				ChecksPollTime:      1 * time.Second,
				ChecksConfig: []config.CheckConfig{
					{
						Type: "icmp",
						CheckProperties: map[string]interface{}{
							"Address": "www.google.com",
							"Timeout": 5 * time.Second,
						},
					},
				},
			}
		})

		It("should serve http requests", func() {
			command := act(checksdConfig)

			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).ShouldNot(HaveOccurred())
			Eventually(session, 10).Should(gbytes.Say("After Check ran"))

			resp, err := http.DefaultClient.Get("http://localhost:8080/")
			Expect(err).ToNot(HaveOccurred())
			contents, err := ioutil.ReadAll(resp.Body)
			Expect(err).ToNot(HaveOccurred())

			Expect(contents).To(MatchJSON(`{}`))
		})

		It("should persis check status to file", func() {
			command := act(checksdConfig)

			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).ShouldNot(HaveOccurred())
			Eventually(session, 10).Should(gbytes.Say("After Check ran"))

			Expect(fileSummaryJsonPath).To(BeAnExistingFile())
			Expect(fileSummaryJsonPath).To(BeARegularFile())

			fileSummaryContents, err := ioutil.ReadFile(fileSummaryJsonPath)
			Expect(err).ToNot(HaveOccurred())
			Expect(fileSummaryContents).To(ContainSubstring(`{"CheckStatus":{"*network.IcmpCheck":[{"Modified":`))
		})

		It("Should run health check", func() {
			command := act(checksdConfig)

			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).ShouldNot(HaveOccurred())
			Eventually(session, 10).Should(gbytes.Say("After Check ran"))
		})

		Context("Check Poll config is not provided", func() {
			BeforeEach(func() {
				checksdConfig.ChecksPollTime = 0
			})

			It("Should default to 30 second value", func() {
				command := act(checksdConfig)

				session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).ShouldNot(HaveOccurred())
				Consistently(session, 25).ShouldNot(gbytes.Say("Before Check ran"))
				Eventually(session, 10).Should(gbytes.Say("Before Check ran"))
			})
		})

		Context("Sending a SIGTERM process signal", func() {
			It("should handle signal and exit gracefully", func() {
				command := act(checksdConfig)

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
			tempConfigFile, err := fileutil.TempFile(os.TempDir(), "checksd", "json")
			Expect(err).ToNot(HaveOccurred())

			err = ioutil.WriteFile(tempConfigFile.Name(), []byte(`not json`), os.ModePerm)
			Expect(err).ToNot(HaveOccurred())

			command = exec.Command(pathToChecksd, fmt.Sprintf("--config=%s", tempConfigFile.Name()))
		})

		It("Should error", func() {
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).ShouldNot(HaveOccurred())
			Eventually(session.Err, 10).Should(gbytes.Say("unable to configure checksd with config file"))
			Eventually(session.ExitCode, 3).Should(Equal(1))
		})
	})
})
