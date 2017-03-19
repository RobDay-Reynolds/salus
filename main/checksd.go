package main

import (
	"encoding/json"
	"github.com/FiloSottile/gvt/fileutils"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"github.com/jessevdk/go-flags"
	"github.com/monkeyherder/moirai/checks"
	"github.com/monkeyherder/moirai/checks/adaptors"
	"github.com/monkeyherder/moirai/checks/writer"
	"github.com/monkeyherder/moirai/config"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const TAG string = "checksd"
const DEFAULT_CHECKS_POLL_TIME time.Duration = 30 * time.Second

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
		exitCode = 1
		return
	}

	config, err := parseConfig(opts)
	if err != nil {
		asyncLog.Error(TAG, "unable to configure checksd with config file", err)
		exitCode = 1
		return
	}

	exitCode = startDaemon(asyncLog, config)
}

func parseConfig(opts *ConfigOpts) (*config.ChecksdConfig, error) {
	configFile, err := os.Open(opts.FilePath)
	if err != nil {
		return nil, err
	}
	configContents, err := ioutil.ReadAll(configFile)
	if err != nil {
		return nil, err
	}

	config := &config.ChecksdConfig{}
	err = json.Unmarshal(configContents, config)
	if err != nil {
		return nil, err
	}

	if config.ChecksPollTime <= 0 {
		config.ChecksPollTime = DEFAULT_CHECKS_POLL_TIME
	}

	return config, nil
}

func startDaemon(logger boshlog.Logger, config *config.ChecksdConfig) int {
	sigChannel := make(chan os.Signal, 8)
	signal.Notify(sigChannel, syscall.SIGTERM, os.Interrupt, os.Kill)

	serverErrChannel := startHealthCheckHttpServerAsync(config.CheckStatusFilePath)
	statusWriter := writer.CheckSummaryWriter{
		PathToCheckSummary: config.CheckStatusFilePath,
		Logger:             logger,
	}
	for {
		time.Sleep(config.ChecksPollTime)
		select {
		case serverErr := <-serverErrChannel:
			logger.Error(TAG, "http server errored with: %v", serverErr)
			return -1
		case sig := <-sigChannel:
			logger.Debug(TAG, "sig received: %v", sig)
			return 0
		default:
			fileutils.RemoveAll(config.CheckStatusFilePath)
			for _, check := range config.Checks {
				checks.Checker(check,
					adaptors.MustPersistCheckStatus(statusWriter, logger),
					adaptors.NewNotifierLogger(logger),
				).Run()
			}
		}
	}
}

func startHealthCheckHttpServerAsync(checkStatusFilePath string) chan error {
	http.HandleFunc("/", http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		contents, _ := ioutil.ReadFile(checkStatusFilePath)
		resp.Write(contents)
	}))

	serverError := make(chan error, 1)
	go func() {
		serverError <- http.ListenAndServe(":8080", nil)
	}()

	return serverError
}
