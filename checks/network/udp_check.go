package network

import (
	"net"
	"fmt"
	"time"
	"github.com/pkg/errors"
)

type UdpCheck struct {
	Port    int
	Timeout time.Duration
}

func (c UdpCheck) Run() error {
	conn, err := net.DialTimeout("udp", fmt.Sprintf("127.0.0.1:%d", c.Port), c.Timeout)
	if err != nil {
		return errors.Wrapf(err, "Port %d is not available", c.Port)
	}

	conn.SetReadDeadline(time.Now().Add(c.Timeout))
	conn.Write([]byte(""))

	_, err = conn.Read(make([]byte, 1))
	if err != nil {
		return errors.Wrapf(err, "Timed out on udp response", c.Port)
	}

	defer conn.Close()

	return nil
}
