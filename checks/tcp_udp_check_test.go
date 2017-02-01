package checks_test

import (
	. "github.com/monkeyherder/moirai/checks"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net"
	"fmt"
	"strconv"
)

var _ = Describe("TcpUdpCheck", func() {

	Describe("TCP", func() {
		var tcpUdpCheck TcpUdpCheck
		var localTcpServer *LocalTcpServer

		BeforeEach(func() {
			localTcpServer = StartLocalTcpServer()

			tcpUdpCheck = TcpUdpCheck{
				Protocol: TCP,
				Port:     localTcpServer.Port,
			}
		})

		AfterEach(func() {
			localTcpServer.Close()
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
})

func StartLocalTcpServer() *LocalTcpServer {
	server := &LocalTcpServer{
		HandleRequest: func(conn net.Conn) {
			conn.Write([]byte("Hello World"))
			conn.Close()
		},
	}
	started := make(chan int, 1)
	go server.Start(0, started)
	<-started

	return server
}

type LocalTcpServer struct {
	Port          int
	HandleRequest func(conn net.Conn)
	listener      net.Listener
}

func (localTcp *LocalTcpServer) Close() error {
	localTcp.HandleRequest = nil
	return localTcp.listener.Close()
}

func (localTcp *LocalTcpServer) Start(port int, started chan int) error {
	defer GinkgoRecover()

	tcpPort := strconv.Itoa(port)

	l, err := net.Listen("tcp", "127.0.0.1:"+tcpPort)
	if err != nil {
		Fail(fmt.Sprintf("Error listening: %v", err.Error()))
	}

	localTcp.listener = l
	_, ephemeralPort, err := net.SplitHostPort(l.Addr().String())
	localTcp.Port, err = strconv.Atoi(ephemeralPort)

	if err != nil {
		Fail(fmt.Sprintf("Cannot determine port: %v", err.Error()))
	}

	fmt.Println("Listening on " + "127.0.0.1:" + ephemeralPort)

	started <- 0
	for {
		conn, err := l.Accept()
		if err == nil && localTcp.HandleRequest != nil {
			localTcp.HandleRequest(conn)
		}
	}
}
