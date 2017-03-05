package network

import (
	"errors"
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"

	"github.com/monkeyherder/moirai/checks"
)

// http://www.iana.org/assignments/protocol-numbers/protocol-numbers.xhtml
const ICMP_PROTOCOL int = 1

type IcmpCheck struct {
	Address string
	Timeout time.Duration
}

func (icmpCheck IcmpCheck) Run() (checks.CheckInfo, error) {
	// Use a non-privileged datagram-oriented ICMP: https://lwn.net/Articles/420800/
	packetConn, err := icmp.ListenPacket("udp4", "0.0.0.0")

	if err != nil {
		return checks.CheckInfo{}, err
	}
	defer packetConn.Close()

	icmpRequest := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID: os.Getpid() & 0xffff, Seq: 1,
			Data: []byte{},
		},
	}
	icmpRequestBody, err := icmpRequest.Marshal(nil)
	if err != nil {
		return checks.CheckInfo{}, err
	}

	ipAddr, err := getAddr(icmpCheck.Address)
	if err != nil {
		return checks.CheckInfo{}, err
	}

	if _, err := packetConn.WriteTo(icmpRequestBody, ipAddr); err != nil {
		return checks.CheckInfo{}, err
	}

	deadline := time.Now().Add(icmpCheck.Timeout)
	packetConn.SetReadDeadline(deadline)

	readBuffer := make([]byte, 1500)
	responseReadLen, _, err := packetConn.ReadFrom(readBuffer)
	if err != nil {
		return checks.CheckInfo{}, err
	}

	icmpResponseMessage, err := icmp.ParseMessage(ICMP_PROTOCOL, readBuffer[:responseReadLen])
	if err != nil {
		return checks.CheckInfo{}, err
	}
	switch icmpResponseMessage.Type {
	case ipv4.ICMPTypeEchoReply:
		return checks.CheckInfo{}, nil
	default:
		return checks.CheckInfo{}, errors.New(fmt.Sprintf("got %+v; expected echo reply", icmpResponseMessage))
	}
}

func getAddr(host string) (net.Addr, error) {
	ips, err := net.LookupIP(host)
	if err != nil {
		return nil, err
	}

	for _, ip := range ips {
		if ip.To4() != nil {
			return &net.UDPAddr{IP: ip}, nil
		}
	}
	return nil, errors.New("no A record")
}
