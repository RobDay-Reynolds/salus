package config_test

import (
	. "github.com/monkeyherder/salus/config"

	"encoding/json"
	"github.com/monkeyherder/salus/checks"
	"github.com/monkeyherder/salus/checks/network"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("ChecksConfig", func() {

	Context("Given a valid checksconfig", func() {
		var unMarshallChecksdConfig *ChecksdConfig
		var checksdConfig *ChecksdConfig
		var checksdConfigJson []byte

		BeforeEach(func() {
			unMarshallChecksdConfig = &ChecksdConfig{}

			checksdConfig = &ChecksdConfig{
				ChecksPollTime: 1 * time.Second,
			}
			var err error
			checksdConfigJson, err = json.Marshal(checksdConfig)
			Expect(err).ToNot(HaveOccurred())
		})

		Context("with network icmp check config", func() {
			BeforeEach(func() {
				checksdConfig.ChecksConfig = []CheckConfig{
					{
						Type: "icmp",
						CheckProperties: map[string]interface{}{
							"Address": "icmp_address",
							"Timeout": 1 * time.Second,
						},
					},
				}
				var err error
				checksdConfigJson, err = json.Marshal(checksdConfig)

				Expect(err).ToNot(HaveOccurred())
			})

			It("should unmarshal valid icmpcheck", func() {
				err := json.Unmarshal(checksdConfigJson, unMarshallChecksdConfig)
				Expect(err).ToNot(HaveOccurred())
				Expect(unMarshallChecksdConfig.Checks).To(HaveLen(1))
				Expect(unMarshallChecksdConfig.Checks[0]).To(Equal(&network.IcmpCheck{
					Address: "icmp_address",
					Timeout: 1 * time.Second,
				}))
			})
		})

		Context("with network socket check config", func() {
			BeforeEach(func() {
				checksdConfig.ChecksConfig = []CheckConfig{
					{
						Type: "unix_socket",
						CheckProperties: map[string]interface{}{
							"Timeout":    1 * time.Second,
							"SocketFile": "unix-socket-file",
						},
					},
				}
				var err error
				checksdConfigJson, err = json.Marshal(checksdConfig)
				Expect(err).ToNot(HaveOccurred())
			})

			It("should unmarshal valid socket check", func() {
				err := json.Unmarshal(checksdConfigJson, unMarshallChecksdConfig)
				Expect(err).ToNot(HaveOccurred())
				Expect(unMarshallChecksdConfig.Checks).To(HaveLen(1))
				Expect(unMarshallChecksdConfig.Checks[0]).To(Equal(&network.UnixSocketCheck{
					Timeout:    1 * time.Second,
					SocketFile: "unix-socket-file",
				}))
			})
		})

		Context("with network tcp check config", func() {
			BeforeEach(func() {
				checksdConfig.ChecksConfig = []CheckConfig{
					{
						Type: "tcp",
						CheckProperties: map[string]interface{}{
							"Timeout": 1 * time.Second,
							"Port":    1234,
						},
					},
				}
				var err error
				checksdConfigJson, err = json.Marshal(checksdConfig)
				Expect(err).ToNot(HaveOccurred())
			})

			It("should unmarshal valid tcp check", func() {
				err := json.Unmarshal(checksdConfigJson, unMarshallChecksdConfig)
				Expect(err).ToNot(HaveOccurred())
				Expect(unMarshallChecksdConfig.Checks).To(HaveLen(1))
				Expect(unMarshallChecksdConfig.Checks[0]).To(Equal(&network.TcpCheck{
					Timeout: 1 * time.Second,
					Port:    1234,
				}))
			})
		})

		Context("with network udp check config", func() {
			BeforeEach(func() {
				checksdConfig.ChecksConfig = []CheckConfig{
					{
						Type: "udp",
						CheckProperties: map[string]interface{}{
							"Timeout": 1 * time.Second,
							"Port":    1234,
						},
					},
				}
				var err error
				checksdConfigJson, err = json.Marshal(checksdConfig)
				Expect(err).ToNot(HaveOccurred())
			})

			It("should unmarshal valid udp check", func() {
				err := json.Unmarshal(checksdConfigJson, unMarshallChecksdConfig)
				Expect(err).ToNot(HaveOccurred())
				Expect(unMarshallChecksdConfig.Checks).To(HaveLen(1))
				Expect(unMarshallChecksdConfig.Checks[0]).To(Equal(&network.UdpCheck{
					Timeout: 1 * time.Second,
					Port:    1234,
				}))
			})
		})

		Context("with file check config", func() {
			BeforeEach(func() {

				checksdConfig.ChecksConfig = []CheckConfig{
					{
						Type: "file",
						CheckProperties: map[string]interface{}{
							"Name":      "name",
							"Path":      "path",
							"IfChanged": "ifchanged",
							"Group":     "group",
							"DependsOn": "dependson",
						},
					},
				}
				var err error
				checksdConfigJson, err = json.Marshal(checksdConfig)
				Expect(err).ToNot(HaveOccurred())
			})

			It("should unmarshal valid file check", func() {
				err := json.Unmarshal(checksdConfigJson, unMarshallChecksdConfig)
				Expect(err).ToNot(HaveOccurred())
				Expect(unMarshallChecksdConfig.Checks).To(HaveLen(1))
				Expect(unMarshallChecksdConfig.Checks[0]).To(Equal(&checks.FileCheck{
					Name:      "name",
					Path:      "path",
					IfChanged: "ifchanged",
					Group:     "group",
					DependsOn: "dependson",
				}))

			})
		})

		Context("with process check config", func() {
			BeforeEach(func() {

				checksdConfig.ChecksConfig = []CheckConfig{
					{
						Type: "process",
						CheckProperties: map[string]interface{}{
							"Name":         "name",
							"Pidfile":      "pidfile",
							"StartProgram": "startprogram",
							"StopProgram":  "stopprogram",
							"Group":        "group",
							"DependsOn":    "dependson",
						},
					},
				}
				var err error
				checksdConfigJson, err = json.Marshal(checksdConfig)
				Expect(err).ToNot(HaveOccurred())
			})

			It("should unmarshal valid process check", func() {
				err := json.Unmarshal(checksdConfigJson, unMarshallChecksdConfig)
				Expect(err).ToNot(HaveOccurred())
				Expect(unMarshallChecksdConfig.Checks).To(HaveLen(1))
				Expect(unMarshallChecksdConfig.Checks[0]).To(Equal(&checks.ProcessCheck{
					Name:         "name",
					Pidfile:      "pidfile",
					StartProgram: "startprogram",
					StopProgram:  "stopprogram",
					Group:        "group",
					DependsOn:    "dependson",
				}))

			})
		})

		Context("multiple different checks", func() {
			BeforeEach(func() {

				checksdConfig.ChecksConfig = []CheckConfig{
					{
						Type: "process",
						CheckProperties: map[string]interface{}{
							"Name":         "name",
							"Pidfile":      "pidfile",
							"StartProgram": "startprogram",
							"StopProgram":  "stopprogram",
							"Group":        "group",
							"DependsOn":    "dependson",
						},
					},
					{
						Type: "icmp",
						CheckProperties: map[string]interface{}{
							"Address": "icmp_address",
							"Timeout": 1 * time.Second,
						},
					},
				}
				var err error
				checksdConfigJson, err = json.Marshal(checksdConfig)
				Expect(err).ToNot(HaveOccurred())
			})

			It("should unmarshal every check", func() {
				err := json.Unmarshal(checksdConfigJson, unMarshallChecksdConfig)
				Expect(err).ToNot(HaveOccurred())
				Expect(unMarshallChecksdConfig.Checks).To(HaveLen(2))
				Expect(unMarshallChecksdConfig.Checks[0]).To(Equal(&checks.ProcessCheck{
					Name:         "name",
					Pidfile:      "pidfile",
					StartProgram: "startprogram",
					StopProgram:  "stopprogram",
					Group:        "group",
					DependsOn:    "dependson",
				}))
				Expect(unMarshallChecksdConfig.Checks[1]).To(Equal(&network.IcmpCheck{
					Address: "icmp_address",
					Timeout: 1 * time.Second,
				}))

			})

		})

		Context("multiple same checks", func() {
			BeforeEach(func() {

				checksdConfig.ChecksConfig = []CheckConfig{
					{
						Type: "icmp",
						CheckProperties: map[string]interface{}{
							"Address": "another_icmp_address",
							"Timeout": 2 * time.Second,
						},
					},
					{
						Type: "icmp",
						CheckProperties: map[string]interface{}{
							"Address": "icmp_address",
							"Timeout": 1 * time.Second,
						},
					},
				}
				var err error
				checksdConfigJson, err = json.Marshal(checksdConfig)
				Expect(err).ToNot(HaveOccurred())
			})

			It("should unmarshal every check", func() {
				err := json.Unmarshal(checksdConfigJson, unMarshallChecksdConfig)
				Expect(err).ToNot(HaveOccurred())
				Expect(unMarshallChecksdConfig.Checks).To(HaveLen(2))
				Expect(unMarshallChecksdConfig.Checks[0]).To(Equal(&network.IcmpCheck{
					Address: "another_icmp_address",
					Timeout: 2 * time.Second,
				}))
				Expect(unMarshallChecksdConfig.Checks[1]).To(Equal(&network.IcmpCheck{
					Address: "icmp_address",
					Timeout: 1 * time.Second,
				}))

			})

		})
	})

	Context("Given a invalid checksconfig", func() {
		var unMarshallChecksdConfig *ChecksdConfig
		var checksdConfigJson []byte

		BeforeEach(func() {
			unMarshallChecksdConfig = &ChecksdConfig{}
			checksdConfigJson = []byte("invalid config")
		})

		It("should return error", func() {
			err := json.Unmarshal(checksdConfigJson, unMarshallChecksdConfig)
			Expect(err).To(HaveOccurred())

		})

		Context("Given multiple checks with some being invalid", func() {
			BeforeEach(func() {
				checksdConfigJson = []byte(`{"checksPollTime":1000000000,
				"checks":[
					{"Type":"icmp", "CheckProperties":{"Address":"icmp_address","Timeout":1000000000}},
					{"Type":"unknown"}
				],"Checks":null}`)
			})

			It("should return an error", func() {
				err := json.Unmarshal(checksdConfigJson, unMarshallChecksdConfig)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Check config with type: 'unknown' is not a valid check"))
			})
		})
	})
})
