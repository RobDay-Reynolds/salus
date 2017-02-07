package network

import (
	"net"
	"fmt"
	"time"
	"github.com/pkg/errors"
)

type UdpCheck struct {
	Port     int
	Timeout  time.Duration
	Protocol UdpProcotol
}

type UdpProcotol func(UdpCon) error

type UdpCon interface {
	Read(b []byte) (n int, err error)
	Write(b []byte) (n int, err error)
	SetDeadline(t time.Time) error
	SetReadDeadline(t time.Time) error
	SetWriteDeadline(t time.Time) error
}

func (c UdpCheck) Run() error {
	conn, err := net.DialTimeout("udp", fmt.Sprintf("127.0.0.1:%d", c.Port), c.Timeout)
	if err != nil {
		return errors.Wrapf(err, "Port %d is not available", c.Port)
	}
	defer conn.Close()

	return c.Protocol(conn)
}

func NewUdpCheck(port int, timeout time.Duration) UdpCheck {
	return UdpCheck{
		Port:    port,
		Timeout: timeout,
		Protocol: func(udpConn UdpCon) error {
			udpConn.SetReadDeadline(time.Now().Add(timeout))
			udpConn.Write([]byte(""))

			_, err := udpConn.Read(make([]byte, 1))
			if err != nil {
				return errors.Wrapf(err, "Timed out on udp response", port)
			}

			return nil
		},
	}
}
