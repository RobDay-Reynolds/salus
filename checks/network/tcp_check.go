package network

import (
	"net"
	"fmt"
	"time"
	"github.com/pkg/errors"
)

type TcpCheck struct {
	Port     int
	Timeout  time.Duration
}

func (c TcpCheck) Run() error {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", c.Port), c.Timeout)
	if err != nil {
		return errors.Wrapf(err, "Port %d is not available", c.Port)
	}

	defer conn.Close()

	return nil
}
