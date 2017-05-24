package emitter_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
	"github.com/nats-io/gnatsd/server"
	"testing"
	"time"
)

func TestNats(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Emitter Suite")
}

// DefaultTestOptions are default options for the unit tests.
var DefaultTestOptions = server.Options{
	Host:           "localhost",
	Port:           4222,
	NoLog:          true,
	NoSigs:         true,
	MaxControlLine: 256,
}

// RunDefaultServer starts a new Go routine based server using the default options
func RunDefaultServer() *server.Server {
	return RunServer(&DefaultTestOptions)
}

// RunServer starts a new Go routine based server
func RunServer(opts *server.Options) *server.Server {
	return RunServerWithAuth(opts, nil)
}

func RunServerWithAuth(opts *server.Options, auth server.Auth) *server.Server {
	if opts == nil {
		opts = &DefaultTestOptions
	}
	s := server.New(opts)
	if s == nil {
		panic("No NATS Server object returned.")
	}

	if auth != nil {
		s.SetClientAuthMethod(auth)
	}

	// Run server in Go routine.
	go s.Start()

	// Wait for accept loop(s) to be started
	if !s.ReadyForConnections(10 * time.Second) {
		panic("Unable to start NATS Server in Go Routine")
	}
	return s
}

// LoadConfig loads a configuration from a filename
func LoadConfig(configFile string) (opts *server.Options) {
	opts, err := server.ProcessConfigFile(configFile)
	if err != nil {
		panic(fmt.Sprintf("Error processing configuration file: %v", err))
	}
	opts.NoSigs, opts.NoLog = true, true
	return
}
