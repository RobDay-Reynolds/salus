package emitter_test

import (
	. "github.com/monkeyherder/salus/emitter"

	"github.com/cloudfoundry/bosh-utils/logger/fakes"
	"github.com/monkeyherder/salus/checks"
	"github.com/monkeyherder/salus/emitter/emitterfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NatsClient", func() {
	var emitter *Emitter
	var fakeNatsClient *emitterfakes.FakeClient

	BeforeEach(func() {
		fakeNatsClient = &emitterfakes.FakeClient{}
		emitter = &Emitter{
			Logger: &fakes.FakeLogger{},
			Client: fakeNatsClient,
		}

		emitter.Start()
	})

	AfterEach(func() {
		emitter.Shutdown()
	})

	It("should emite checks to nats client", func() {
		checkChannel := emitter.EmitCheck()

		checkChannel <- checks.CheckInfo{
			Note:   "test note",
			Status: "test status",
		}

		Eventually(func() int {
			return fakeNatsClient.PublishCallCount()
		}).Should(Equal(1))

		subject, payload := fakeNatsClient.PublishArgsForCall(0)
		Expect(subject).To(Equal("hm.salus.check.uuid"))
		Expect(payload).To(MatchJSON(`{
							"Status": "test status",
							"Note"  : "test note"
						   }`))
	})
})
