package network_test

import (
	. "github.com/monkeyherder/moirai/checks/network"

	"fmt"
	"net"
	"strconv"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UdpCheck", func() {
	var udpCheck UdpCheck
	var localUdpServer *LocalUdpServer

	BeforeEach(func() {
		localUdpServer = StartLocalUdpServer()

		udpCheck = NewUdpCheck(localUdpServer.Port, 2*time.Second)
	})

	AfterEach(func() {
		localUdpServer.CloseUdp()
	})

	Context("A port that has a udp service listening on it", func() {
		Context("And is rechable and responsive", func() {
			BeforeEach(func() {
				localUdpServer.ShouldRespond = true
			})

			It("Check should return as healthy", func() {
				_, _, err := udpCheck.Run()
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("And is not responsive", func() {
			BeforeEach(func() {
				localUdpServer.ShouldRespond = false
			})
			It("Check should return as not healthy", func() {
				_, _, err := udpCheck.Run()
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Context("A Port that is not a healthy udp", func() {
		BeforeEach(func() {
			udpCheck.Port = 1
		})

		It("Check should return as unhealthy", func() {
			_, _, err := udpCheck.Run()
			Expect(err).To(HaveOccurred())
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
	ShouldRespond  bool
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

	started <- 0
	for {
		_, addr, _ := l.ReadFrom(make([]byte, 1))
		if localTcpUdp.ShouldRespond {
			_, err = l.WriteTo([]byte("foo"), addr)
		}
	}
}
