package monit_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/monkeyherder/moirai/monit"
)

var _ = Describe("ReadMonitFile", func() {
	It("creates MonitFile struct with Check structs for a check in the file", func() {
		monitFile, err := ReadMonitFile(fixturesPath + "/simple.monit")

		Expect(err).To(BeNil())

		simpleCheck := ProcessCheck{
			Pidfile:      "/path/to/pid",
			StartProgram: "/path/to/start/command",
			StopProgram:  "/path/to/command with args",
			Group:        "test_group",
		}

		Expect(monitFile.Checks).To(Equal([]Check{simpleCheck}))
	})
})
