package adaptors_test

import (
	. "github.com/monkeyherder/moirai/checks/adaptors"

	"encoding/json"
	"errors"
	"github.com/monkeyherder/moirai/checks"
	"github.com/monkeyherder/moirai/checks/checksfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
)

var _ = Describe("PersistCheckStatusSummary", func() {

	Context("Given a check", func() {
		var persistCheckAdaptor checks.CheckAdaptor
		var fakeCheck *checksfakes.FakeCheck
		var pathToCheckSummaryFile string

		BeforeEach(func() {
			tempFile, err := ioutil.TempFile(os.TempDir(), "checksummarytestfile")
			Expect(err).ToNot(HaveOccurred())
			pathToCheckSummaryFile = tempFile.Name()

			persistCheckAdaptor = PersistCheckSummary(pathToCheckSummaryFile)
			fakeCheck = &checksfakes.FakeCheck{}
		})

		Context("with a failing status response", func() {
			BeforeEach(func() {
				failCheckInfo := checks.CheckInfo{
					Status: "Not good",
					Note:   "didn't go well",
				}
				fakeCheck.RunReturns(failCheckInfo, errors.New("error too"))
			})

			It("Should persist the check response with a failing status", func() {
				checkInfo, checkErr := persistCheckAdaptor(fakeCheck).Run()

				Expect(checkInfo.Status).To(Equal("Not good"))
				Expect(checkInfo.Note).To(Equal("didn't go well"))
				Expect(checkErr).To(HaveOccurred())
				Expect(checkErr.Error()).To(ContainSubstring("error too"))

				Expect(pathToCheckSummaryFile).To(BeAnExistingFile())
				checkSummaryContents, err := ioutil.ReadFile(pathToCheckSummaryFile)
				Expect(err).ToNot(HaveOccurred())

				actualCheckSummary := &CheckSummary{}
				json.Unmarshal(checkSummaryContents, actualCheckSummary)

				//Expect(actualCheckSummary.CheckSummary).To(HaveLen(1))
			})
		})

	})
})
