package checks

import (
	"net"
	"fmt"
	"time"
	"github.com/pkg/errors"
)

const (
	TCP int = iota
	UDP
)

type TcpUdpCheck struct {
	Protocol int
	Port     int
	Timeout  time.Duration
}

func (c TcpUdpCheck) Run() error {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", c.Port), c.Timeout)
	if err != nil {
		return errors.Wrapf(err, "Port %d is not available", c.Port)
	}

	defer conn.Close()

	return nil
}
