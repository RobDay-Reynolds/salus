package network

import (
	"github.com/pkg/errors"
	"net"
	"time"
)

type UnixSocketCheck struct {
	SocketFile string
	Timeout    time.Duration
}

func (u UnixSocketCheck) Run() error {
	c, err := net.Dial("unix", u.SocketFile)
	if err != nil {
		return errors.Wrap(err, "Unable to connect to unix socket file")
	}
	defer c.Close()

	return nil
}
