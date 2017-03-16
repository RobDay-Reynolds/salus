package writer

import (
	"encoding/json"
	"github.com/cloudfoundry/bosh-utils/logger"
	"github.com/golang/go/src/pkg/io/ioutil"
	"github.com/monkeyherder/moirai/checks/adaptors"
	"os"
)

const TAG string = "CheckSummaryWriter"

type CheckStatus struct {
	CheckStatus map[string][]adaptors.Status
}

type CheckSummaryWriter struct {
	PathToCheckSummary string
	Logger             logger.Logger
}

func (csw CheckSummaryWriter) Write(status adaptors.Status) error {
	checkSummaryFile, err := os.OpenFile(csw.PathToCheckSummary, os.O_RDWR|os.O_CREATE, 0600)

	fileContents, err := ioutil.ReadAll(checkSummaryFile)
	if err != nil {
		csw.Logger.Error(TAG, "unable to read status summary file: %s", err.Error())
		return err
	}

	var checkStatus CheckStatus
	err = json.Unmarshal(fileContents, &checkStatus)
	if err != nil {
		csw.Logger.Debug(TAG, "unable to unmarshal json into summary file (%s). Initializing a new one: %s", fileContents, err.Error())

		checkStatus = CheckStatus{
			CheckStatus: map[string][]adaptors.Status{},
		}
	}

	checkStatus.CheckStatus[status.CheckType] = append(checkStatus.CheckStatus[status.CheckType], status)

	checkStatusJson, err := json.Marshal(checkStatus)
	if err != nil {
		csw.Logger.Error(TAG, "unable to marshal check status into json: %s", err.Error())
		return err
	}

	tempFile, err := ioutil.TempFile(os.TempDir(), "intermediate_summary_file")
	if err != nil {
		csw.Logger.Error(TAG, "unable to create temp file: %s", err.Error())
		return err
	}

	_, err = tempFile.Write(checkStatusJson)
	if err != nil {
		csw.Logger.Error(TAG, "unable to write json into temp file: %s", err.Error())
		return err
	}

	return os.Rename(tempFile.Name(), checkSummaryFile.Name())
}
