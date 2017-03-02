package network

import (
	"fmt"
	"net"
	"time"

	"github.com/pkg/errors"

	"github.com/monkeyherder/moirai/checks"
)

type TcpCheck struct {
	Port    int
	Timeout time.Duration
}

func (c TcpCheck) Run() (checks.CheckInfo, error) {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", c.Port), c.Timeout)
	if err != nil {
		return checks.CheckInfo{}, errors.Wrapf(err, "Port %d is not available", c.Port)
	}

	defer conn.Close()

	return checks.CheckInfo{}, nil
}
