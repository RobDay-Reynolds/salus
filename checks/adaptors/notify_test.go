package adaptors_test

import (
	. "github.com/monkeyherder/moirai/checks/adaptors"

	"errors"

	"github.com/monkeyherder/moirai/checks"
	"github.com/monkeyherder/moirai/checks/adaptors/adaptorsfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Notify", func() {
	var notifier checks.CheckAdaptor
	var fakeNotifyHook *adaptorsfakes.FakeNotifier

	BeforeEach(func() {
		fakeNotifyHook = &adaptorsfakes.FakeNotifier{}
		notifier = Notify(fakeNotifyHook)
	})

	Context("Given a check", func() {
		var successCheckFn checks.Check
		var checkWithNotifier checks.Check

		BeforeEach(func() {
			successCheckFn = checks.CheckFunc(func() (checks.CheckInfo, error) {
				return checks.CheckInfo{}, nil
			})
			checkWithNotifier = notifier(successCheckFn)
		})

		It("should notify that the check has run", func() {
			checkWithNotifier.Run()

			Expect(fakeNotifyHook.BeforeCheckCallCount()).To(Equal(1))
			checkArgsForCall := fakeNotifyHook.BeforeCheckArgsForCall(0)
			Expect(checkArgsForCall).To(BeAssignableToTypeOf(successCheckFn))

			Expect(fakeNotifyHook.OnErrorCallCount()).To(Equal(0))

			Expect(fakeNotifyHook.AfterCheckCallCount()).To(Equal(1))
			checkArgsForCall = fakeNotifyHook.AfterCheckArgsForCall(0)
			Expect(checkArgsForCall).To(BeAssignableToTypeOf(successCheckFn))
		})

		Context("A failing check", func() {

			BeforeEach(func() {
				successCheckFn := checks.CheckFunc(func() (checks.CheckInfo, error) {
					return checks.CheckInfo{}, errors.New("some error")
				})
				checkWithNotifier = notifier(successCheckFn)
			})

			It("Should notify that the check failed", func() {
				checkWithNotifier.Run()

				Expect(fakeNotifyHook.BeforeCheckCallCount()).To(Equal(1))

				Expect(fakeNotifyHook.OnErrorCallCount()).To(Equal(1))
				checkArgsForCall, errArgForCall := fakeNotifyHook.OnErrorArgsForCall(0)
				Expect(checkArgsForCall).To(BeAssignableToTypeOf(successCheckFn))
				Expect(errArgForCall).To(HaveOccurred())
				Expect(errArgForCall).To(MatchError(errors.New("some error")))

				Expect(fakeNotifyHook.AfterCheckCallCount()).To(Equal(1))
			})
		})
	})

})
