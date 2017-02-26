package main

import (
	"encoding/json"
	"github.com/golang/go/src/pkg/io/ioutil"
	"github.com/jessevdk/go-flags"
	"github.com/monkeyherder/moirai/checks"
	"github.com/monkeyherder/moirai/checks/adaptors"
	"github.com/monkeyherder/moirai/checks/network"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const TAG string = "checksd"

type ChecksdConfig struct {
	ChecksPollTime time.Duration `json:"checksPollTime"`
	IcmpChecks     []network.IcmpCheck `json:"icmpChecks"`
}

type ConfigOpts struct {
	FilePath string `short:"c" long:"config" description:"path to checksd config" value-name:"FILE"`
}

func main() {
	exitCode := 0
	asyncLog := boshlog.NewAsyncWriterLogger(boshlog.LevelDebug, os.Stdout, os.Stderr)
	defer func() {
		asyncLog.FlushTimeout(time.Minute)
		os.Exit(exitCode)
	}()

	opts := &ConfigOpts{}
	_, err := flags.ParseArgs(opts, os.Args[1:])
	if err != nil {
		return
	}

	config, err := parseConfig(opts)
	if err != nil {
		asyncLog.Error(TAG, "unable to configure checksd with config file: ", err.Error())
		return
	}

	exitCode = startDaemon(asyncLog, ChecksdConfig{
		IcmpChecks: config.IcmpChecks,
	})
}

func parseConfig(opts *ConfigOpts) (*ChecksdConfig, error) {
	configFile, err := os.Open(opts.FilePath)
	if err != nil {
		return nil, err
	}
	configContents, err := ioutil.ReadAll(configFile)
	if err != nil {
		return nil, err
	}

	config := &ChecksdConfig{}
	err = json.Unmarshal(configContents, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func startDaemon(logger boshlog.Logger, config ChecksdConfig) int {
	sigChannel := make(chan os.Signal, 8)
	signal.Notify(sigChannel, syscall.SIGTERM, os.Interrupt, os.Kill)

	for {
		time.Sleep(1 * time.Second)
		select {
		case sig := <-sigChannel:
			logger.Debug(TAG, "sig received: %v", sig)
			return 0
		default:
			for _, icmpCheck := range config.IcmpChecks {
				checks.Checker(icmpCheck, adaptors.NewNotifierLogger(logger)).Run()
			}
		}
	}
}
