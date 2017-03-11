package adaptors

import (
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"github.com/monkeyherder/moirai/checks"
)

const TAG string = "PersistCheckStatus"

type Status struct {
	CheckInfo  checks.CheckInfo
	CheckError error
}

type CheckStatus struct {
	CheckStatus []Status
}

type CheckStatusWriter interface {
	Write(summary Status) error
}

func PersistCheckStatus(checkSummaryWriter CheckStatusWriter, Logger boshlog.Logger) checks.CheckAdaptor {
	return func(check checks.Check) checks.Check {
		return checks.CheckFunc(func() (checks.CheckInfo, error) {
			checkInfo, err := check.Run()
			writerErr := checkSummaryWriter.Write(Status{
				CheckInfo:  checkInfo,
				CheckError: err,
			})

			if writerErr != nil {
				Logger.Error(TAG, "Unable to write check status. Exiting...", writerErr)
				panic("Unable to persist checkinfo")
			}

			return checkInfo, err
		})
	}
}
