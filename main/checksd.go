package main

import (
	"github.com/monkeyherder/moirai/checks/network"
	"github.com/monkeyherder/moirai/checks"
	"time"
	"github.com/monkeyherder/moirai/checks/adaptors"
	"os"
	"fmt"
	"syscall"
	"os/signal"
)

func main() {
	icmpCheck := network.IcmpCheck{
		Address: "www.google.com",
		Timeout: 5 * time.Second,
	}

	sigChannel := make(chan os.Signal, 8)
	signal.Notify(sigChannel, syscall.SIGTERM, os.Interrupt, os.Kill)

	for {
		time.Sleep(1 * time.Second)
		select {
		case sig := <-sigChannel:
			fmt.Println(fmt.Sprintf("sig received: %v", sig))
			os.Exit(0)
		default:
			checks.Checker(icmpCheck, adaptors.NewNotifierLogger()).Run()
		}
	}
}
