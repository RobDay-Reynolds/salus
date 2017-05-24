package emitter_test

import (
	"github.com/monkeyherder/salus/emitter"
	"github.com/nats-io/gnatsd/server"
	"github.com/nats-io/go-nats"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {
	var natsClient emitter.NatsClient
	var natsServer *server.Server

	Context("Given a healthy nats server", func() {
		BeforeEach(func() {
			natsServer = RunDefaultServer()
			var err error
			natsClient, err = emitter.NewNatsClient(nats.DefaultURL)
			Expect(err).ToNot(HaveOccurred())
		})

		AfterEach(func() {
			natsServer.Shutdown()
		})

		It("should be able to connect to nats server", func() {
			err := natsClient.Publish("subject", []byte("test"))
			Expect(err).ToNot(HaveOccurred())

			Expect(natsServer.NumClients()).To(Equal(1))
		})
	})
})
