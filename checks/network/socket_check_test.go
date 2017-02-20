package network_test

import (
	. "github.com/monkeyherder/moirai/checks/network"

	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net"
	"os"
	"path/filepath"
)

var _ = Describe("SocketCheck", func() {

	var unixSocketCheck UnixSocketCheck
	var localSocketServer *LocalSocketServer

	BeforeEach(func() {
		localSocketServer = StartLocalUnixSocketServer()

		unixSocketCheck = UnixSocketCheck{
			SocketFile: localSocketServer.Socket,
		}
	})

	AfterEach(func() {
		localSocketServer.Close()
	})

	Context("A Socket that is accessible and responsive", func() {
		It("Check should return as healthy", func() {
			err := unixSocketCheck.Run()
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("A Socket that cannot be accessed", func() {
		BeforeEach(func() {
			unixSocketCheck.SocketFile = "/a-non-existent-socket.sock"
		})

		It("Check should return as unhealthy", func() {
			err := unixSocketCheck.Run()
			Expect(err).To(HaveOccurred())
		})
	})
})

func StartLocalUnixSocketServer() *LocalSocketServer {
	server := &LocalSocketServer{
		HandleRequest: echoServer,
	}
	started := make(chan int, 1)
	go server.Start(started)
	<-started

	return server
}

type LocalSocketServer struct {
	Listener      net.Listener
	Socket        string
	HandleRequest func(conn net.Conn)
}

func (localSocketServer *LocalSocketServer) Close() error {
	localSocketServer.HandleRequest = nil
	return localSocketServer.Listener.Close()
}

func (localSocketServer *LocalSocketServer) Start(started chan int) {
	defer GinkgoRecover()
	defer func() {
		started <- -1
	}()

	socketFileName := filepath.Join(os.TempDir(), "echo.sock")

	localSocketServer.Socket = socketFileName
	l, err := net.Listen("unix", socketFileName)
	if err != nil {
		Expect(err).ToNot(HaveOccurred())
	}
	localSocketServer.Listener = l

	started <- 0
	for {
		fd, _ := l.Accept()
		if localSocketServer.HandleRequest != nil {
			localSocketServer.HandleRequest(fd)
		}
	}
}

func echoServer(c net.Conn) {
	defer GinkgoRecover()

	for {
		buf := make([]byte, 512)
		nr, err := c.Read(buf)
		if err != nil {
			return
		}

		data := buf[0:nr]
		fmt.Printf("Received: %v", string(data))
		_, err = c.Write(data)
		if err != nil {
			Expect(err).ToNot(HaveOccurred())
		}
	}
}
