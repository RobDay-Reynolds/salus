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

			simpleCheck := ProcessCheck{
				Name:         "test_process",
				Pidfile:      "/path/to/test/pid",
				StartProgram: "/path/to/test/start/command",
				StopProgram:  "/path/to/test/command with args",
				Group:        "test_group",
			}

			anotherCheck := FileCheck{
				Name:      "file_check",
				Path:      "/path/to/file",
				IfChanged: "/path/to/command",
				Group:     "test_group",
			}

			Expect(monitFile.Checks[0]).To(Equal(simpleCheck))
			Expect(monitFile.Checks[1]).To(Equal(anotherCheck))
		})
	})

	Context("with a re-ordered monit file", func() {
		It("creates MonitFile struct with Check struct for a check in the file", func() {
			monitFile, err := ReadMonitFile(fixturesPath + "/unordered.monit")

			Expect(err).To(BeNil())

			simpleCheck := ProcessCheck{
				Name:         "test_process",
				Pidfile:      "/path/to/test/pid",
				StartProgram: "/path/to/test/start/command",
				StopProgram:  "/path/to/test/command with args",
				Group:        "test_group",
			}

			Expect(monitFile.Checks[0]).To(Equal(simpleCheck))
		})
	})
})
