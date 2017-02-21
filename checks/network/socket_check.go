package network

import (
	"net"
	"time"

	"github.com/pkg/errors"
)

type UnixSocketCheck struct {
	SocketFile string
	Timeout    time.Duration
}

func (u UnixSocketCheck) Run() (string, string, error) {
	c, err := net.Dial("unix", u.SocketFile)
	if err != nil {
		return "", "", errors.Wrap(err, "Unable to connect to unix socket file")
	}
	defer c.Close()

	return "", "", nil
}
