package main

import (
	"fmt"
	"github.com/monkeyherder/moirai/checks"
	"github.com/monkeyherder/moirai/checks/adaptors"
	"github.com/monkeyherder/moirai/checks/network"
	"os"
	"os/signal"
	"syscall"
	"time"
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
