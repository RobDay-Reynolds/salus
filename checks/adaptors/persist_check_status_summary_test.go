package adaptors_test

import (
	. "github.com/monkeyherder/moirai/checks/adaptors"

	"errors"
	"github.com/cloudfoundry/bosh-utils/logger/loggerfakes"
	"github.com/monkeyherder/moirai/checks"
	"github.com/monkeyherder/moirai/checks/adaptors/adaptorsfakes"
	"github.com/monkeyherder/moirai/checks/checksfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PersistCheckStatusSummary", func() {

	Context("Given a check", func() {
		var persistCheckAdaptor checks.CheckAdaptor
		var fakeCheck *checksfakes.FakeCheck
		var fakeCheckStatusWriter *adaptorsfakes.FakeCheckStatusWriter

		BeforeEach(func() {
			fakeCheckStatusWriter = &adaptorsfakes.FakeCheckStatusWriter{}

			persistCheckAdaptor = PersistCheckStatus(fakeCheckStatusWriter, &loggerfakes.FakeLogger{})
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

				checkStatusArg := fakeCheckStatusWriter.WriteArgsForCall(0)
				Expect(checkStatusArg.CheckInfo.Status).To(Equal("Not good"))
				Expect(checkStatusArg.CheckInfo.Note).To(Equal("didn't go well"))
				Expect(checkStatusArg.CheckError).To(HaveOccurred())
				Expect(checkStatusArg.CheckError.Error()).To(ContainSubstring("error too"))
			})

			Context("with the check status writer failing to write", func() {
				BeforeEach(func() {
					fakeCheckStatusWriter.WriteReturns(errors.New("unable to write the check summary due to reasons"))
				})

				It("should panic", func() {
					Expect(func() { persistCheckAdaptor(fakeCheck).Run() }).To(Panic())
				})
			})
		})

		Context("with a healthy status response", func() {
			BeforeEach(func() {
				failCheckInfo := checks.CheckInfo{
					Status: "All good",
					Note:   "went well",
				}
				fakeCheck.RunReturns(failCheckInfo, nil)
			})

			It("should not write checkinfo with an error", func() {
				persistCheckAdaptor(fakeCheck).Run()

				checkStatusArg := fakeCheckStatusWriter.WriteArgsForCall(0)
				Expect(checkStatusArg.CheckInfo.Status).To(Equal("All good"))
				Expect(checkStatusArg.CheckInfo.Note).To(Equal("went well"))
				Expect(checkStatusArg.CheckError).ToNot(HaveOccurred())
			})
		})

	})
})
