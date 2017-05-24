package emitter

import "github.com/nats-io/go-nats"

type NatsClient struct {
	nc *nats.Conn
}

func NewNatsClient(url string) (NatsClient, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return NatsClient{}, err
	}

	return NatsClient{nc: nc}, nil
}

func (n NatsClient) Publish(subject string, bytes []byte) error {
	return n.nc.Publish(subject, bytes)
}
