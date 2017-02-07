package network_test

import (
	. "github.com/monkeyherder/moirai/checks/network"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net"
	"fmt"
	"strconv"
)

var _ = Describe("TcpCheck", func() {
	var tcpUdpCheck TcpCheck
	var localTcpServer *LocalTcpUdpServer

	BeforeEach(func() {
		localTcpServer = StartLocalTcpServer()

		tcpUdpCheck = TcpCheck{
			Port:     localTcpServer.Port,
		}
	})

	AfterEach(func() {
		localTcpServer.CloseTcp()
	})

	Context("A Port that is rechable and responsive", func() {
		It("Check should return as healthy", func() {
			err := tcpUdpCheck.Run()
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("A Port that is not rechable", func() {
		BeforeEach(func() {
			tcpUdpCheck.Port = 1
		})

		It("Check should return as unhealthy", func() {
			err := tcpUdpCheck.Run()
			Expect(err).To(HaveOccurred())
		})
	})

})

func StartLocalTcpServer() *LocalTcpUdpServer {
	server := &LocalTcpUdpServer{
		Protocol: "tcp",
		HandleRequest: func(conn net.Conn) {
			conn.Write([]byte("Hello World"))
			conn.Close()
		},
	}
	started := make(chan int, 1)
	go server.StartTcp(0, started)
	<-started

	return server
}

type LocalTcpUdpServer struct {
	Port           int
	Protocol       string
	HandleRequest  func(conn net.Conn)
	listener       net.Listener
	listenerPacket net.PacketConn
}

func (localTcpUdp *LocalTcpUdpServer) CloseTcp() error {
	localTcpUdp.HandleRequest = nil
	return localTcpUdp.listener.Close()
}

func (localTcpUdp *LocalTcpUdpServer) StartTcp(port int, started chan int) error {
	defer GinkgoRecover()

	tcpPort := strconv.Itoa(port)

	l, err := net.Listen("tcp", "127.0.0.1:"+tcpPort)
	if err != nil {
		Fail(fmt.Sprintf("Error listening: %v", err.Error()))
	}

	localTcpUdp.listener = l
	_, ephemeralPort, err := net.SplitHostPort(l.Addr().String())
	localTcpUdp.Port, err = strconv.Atoi(ephemeralPort)

	if err != nil {
		Fail(fmt.Sprintf("Cannot determine port: %v", err.Error()))
	}

	fmt.Println("Listening on " + "127.0.0.1:" + ephemeralPort)

	started <- 0
	for {
		conn, err := l.Accept()
		if err == nil && localTcpUdp.HandleRequest != nil {
			localTcpUdp.HandleRequest(conn)
		}
	}
}
