package checks_test

import (
	"io/ioutil"
	"syscall"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/monkeyherder/moirai/checks"
)

var _ = Describe("FileCheck", func() {
	var check FileCheck

	Describe("Run", func() {
		Context("when specified file exists", func() {
			It("returns nil", func() {
				file, err := ioutil.TempFile("", "testmainconf")
				Expect(err).ToNot(HaveOccurred())
				defer syscall.Unlink(file.Name())

				check = FileCheck{
					Name: "TestCheck",
					Path: file.Name(),
				}

				err = check.Run()
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when specified file does not exist", func() {
			It("returns an error", func() {
				check = FileCheck{
					Name: "TestCheck",
					Path: "/path/to/unknown/file",
				}

				err := check.Run()
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
