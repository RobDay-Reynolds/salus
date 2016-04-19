package monit_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/monkeyherder/moirai/monit"
)

var _ = Describe("ReadMonitFile", func() {
	Context("with a normally formatted monit file", func() {
		It("creates MonitFile struct with Check structs for each check in the file", func() {
			monitFile, err := ReadMonitFile(fixturesPath + "/simple.monit")

			Expect(err).To(BeNil())

			fileCheck := FileCheck{
				Name:      "file_check",
				Path:      "/path/to/file",
				IfChanged: "/path/to/command",
				Group:     "test_group",
				DependsOn: "something_else",
			}

			failedSocket := FailedSocket{
				SocketFile: "/path/to/socket.sock",
				Timeout:    5,
				NumCycles:  5,
				Action:     "restart",
			}

			failedHost := FailedHost{
				Host:      "1.2.3.4",
				Port:      "9876",
				Protocol:  "http",
				Timeout:   20,
				NumCycles: 10,
				Action:    "stop",
			}

			processCheck := ProcessCheck{
				Name:         "test_process",
				Pidfile:      "/path/to/test/pid",
				StartProgram: "/path/to/test/start/command",
				StopProgram:  "/path/to/test/command with args",
				FailedSocket: failedSocket,
				FailedHost:   failedHost,
				Group:        "test_group",
				DependsOn:    "file_check",
			}

			Expect(monitFile.Checks[0]).To(Equal(fileCheck))
			Expect(monitFile.Checks[1]).To(Equal(processCheck))
		})
	})

	Context("with a re-ordered monit file", func() {
		It("creates MonitFile struct with Check struct for a check in the file", func() {
			monitFile, err := ReadMonitFile(fixturesPath + "/unordered.monit")

			Expect(err).To(BeNil())

			failedSocket := FailedSocket{
				SocketFile: "/path/to/another/socket.sock",
				Timeout:    60,
				NumCycles:  15,
				Action:     "stop",
			}

			simpleCheck := ProcessCheck{
				Name:         "test_process",
				Pidfile:      "/path/to/test/pid",
				StartProgram: "/path/to/test/start/command",
				StopProgram:  "/path/to/test/command with args",
				FailedSocket: failedSocket,
				Group:        "test_group",
			}

			Expect(monitFile.Checks[0]).To(Equal(simpleCheck))
		})
	})

	Context("with non-standard check lengths", func() {
		It("creates MonitFile struct with Check struct for all checks in the file", func() {
			monitFile, err := ReadMonitFile(fixturesPath + "/short_entries.monit")

			Expect(err).To(BeNil())

			shortCheck := ProcessCheck{
				Name:    "short_process",
				Pidfile: "/path/to/short/pid",
			}

			anotherCheck := ProcessCheck{
				Name:         "another_process",
				Pidfile:      "/path/to/another/pid",
				StartProgram: "/path/to/short/start/command",
			}

			Expect(monitFile.Checks[0]).To(Equal(shortCheck))
			Expect(monitFile.Checks[1]).To(Equal(anotherCheck))
		})
	})
})
