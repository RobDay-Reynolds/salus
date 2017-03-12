package writer

import (
	"encoding/json"
	"github.com/golang/go/src/pkg/io/ioutil"
	"github.com/monkeyherder/moirai/checks/adaptors"
	"os"
)

type CheckStatus struct {
	CheckStatus map[string][]adaptors.Status
}

type CheckSummaryWriter struct {
	PathToCheckSummary string
}

func (csw CheckSummaryWriter) Write(status adaptors.Status) error {
	file, err := os.OpenFile(csw.PathToCheckSummary, os.O_RDWR|os.O_CREATE, 0600)

	fileContents, err := ioutil.ReadAll(file)
	if err != nil {

	}

	var checkStatus *CheckStatus
	err = json.Unmarshal(fileContents, checkStatus)
	if err != nil {
		checkStatus = &CheckStatus{
			CheckStatus: map[string][]adaptors.Status{status.CheckType: {status}},
		}
	} else {
		checkStatus.CheckStatus[status.CheckType] = append(checkStatus.CheckStatus[status.CheckType], status)
	}

	checkStatusJson, err := json.Marshal(checkStatus)
	if err != nil {
		return err
	}
	_, err = file.Write(checkStatusJson)
	return err
}
