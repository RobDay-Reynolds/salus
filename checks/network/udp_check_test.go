package network_test

import (
	. "github.com/monkeyherder/moirai/checks/network"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
	"net"
	"strconv"
	"fmt"
)

var _ = Describe("UdpCheck", func() {
	Describe("UDP", func() {
		var udpCheck UdpCheck
		var localUdpServer *LocalUdpServer

		BeforeEach(func() {
			localUdpServer = StartLocalUdpServer()

			udpCheck = UdpCheck{
				Port:     localUdpServer.Port,
				Timeout:  2 * time.Second,
			}
		})

		AfterEach(func() {
			localUdpServer.CloseUdp()
		})

		Context("A Port that is rechable and responsive", func() {
			It("Check should return as healthy", func() {
				err := udpCheck.Run()
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("A Port that is not a healthy udp", func() {
			BeforeEach(func() {
				udpCheck.Port = 1
			})

			It("Check should return as unhealthy", func() {
				err := udpCheck.Run()
				Expect(err).To(HaveOccurred())
			})
		})

	})
})

func StartLocalUdpServer() *LocalUdpServer {
	server := &LocalUdpServer{
		Protocol: "udp",
	}
	started := make(chan int, 1)
	go server.StartUdp(0, started)
	<-started

	return server
}

type LocalUdpServer struct {
	Port           int
	Protocol       string
	listenerPacket net.PacketConn
}

func (localTcpUdp *LocalUdpServer) CloseUdp() error {
	return localTcpUdp.listenerPacket.Close()
}

func (localTcpUdp *LocalUdpServer) StartUdp(port int, started chan int) error {
	defer GinkgoRecover()

	portToListenOn := strconv.Itoa(port)

	resAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:"+portToListenOn)
	l, err := net.ListenUDP("udp", resAddr)

	if err != nil {
		Fail(fmt.Sprintf("Error listening: %v", err.Error()))
	}

	localTcpUdp.listenerPacket = l
	_, ephemeralPort, err := net.SplitHostPort(l.LocalAddr().String())
	localTcpUdp.Port, err = strconv.Atoi(ephemeralPort)

	if err != nil {
		Fail(fmt.Sprintf("Cannot determine port: %v", err.Error()))
	}

	fmt.Println("Listening on " + "127.0.0.1:" + ephemeralPort)

	started <- 0
	for {
		_, addr, _ := l.ReadFrom(make([]byte, 1))
		_, err = l.WriteTo([]byte("foo"), addr)
	}

	fmt.Println(err)

	return nil
}
