package network

import (
	"net"
	"time"

	"github.com/pkg/errors"

	"github.com/monkeyherder/moirai/checks"
)

type UnixSocketCheck struct {
	SocketFile string
	Timeout    time.Duration
}

func (u UnixSocketCheck) Run() (checks.CheckInfo, error) {
	c, err := net.Dial("unix", u.SocketFile)
	if err != nil {
		return checks.CheckInfo{}, errors.Wrap(err, "Unable to connect to unix socket file")
	}
	defer c.Close()

	return checks.CheckInfo{}, nil
}
