package network

import (
	"fmt"
	"github.com/pkg/errors"
	"net"
	"time"
)

type TcpCheck struct {
	Port    int
	Timeout time.Duration
}

func (c TcpCheck) Run() error {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", c.Port), c.Timeout)
	if err != nil {
		return errors.Wrapf(err, "Port %d is not available", c.Port)
	}

	defer conn.Close()

	return nil
}
