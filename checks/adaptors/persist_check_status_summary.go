package adaptors

import (
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"github.com/monkeyherder/moirai/checks"
	"reflect"
	"time"
)

const TAG string = "PersistCheckStatus"

type Status struct {
	Modified   time.Time
	CheckType  string `json:"-"`
	CheckInfo  checks.CheckInfo
	CheckError string
}

type CheckStatusWriter interface {
	Write(status Status) error
}

func MustPersistCheckStatus(checkSummaryWriter CheckStatusWriter, Logger boshlog.Logger) checks.CheckAdaptor {
	return func(check checks.Check) checks.Check {
		return checks.CheckFunc(func() (checks.CheckInfo, error) {
			checkInfo, err := check.Run()

			var errString string
			if err != nil {
				errString = err.Error()
			}
			writerErr := checkSummaryWriter.Write(Status{
				CheckType:  reflect.TypeOf(check).String(),
				Modified:   time.Now(),
				CheckInfo:  checkInfo,
				CheckError: errString,
			})

			if writerErr != nil {
				Logger.Error(TAG, "Unable to write check status. Exiting...", writerErr)
				panic("Unable to persist checkinfo")
			}

			return checkInfo, err
		})
	}
}
