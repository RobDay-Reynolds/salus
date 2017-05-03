package adaptors

import (
	"reflect"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"github.com/monkeyherder/salus/checks"
)

func NewNotifierLogger(logger boshlog.Logger) checks.CheckAdaptor {
	return Notify(notifierLogger{
		Logger: logger,
	})
}

type Notifier interface {
	BeforeCheck(checks.Check)
	AfterCheck(checks.Check)
	OnError(checks.Check, error)
}

func Notify(notifier Notifier) checks.CheckAdaptor {
	return func(check checks.Check) checks.Check {
		return checks.CheckFunc(func() (checks.CheckInfo, error) {
			notifier.BeforeCheck(check)
			checkInfo, err := check.Run()
			if err != nil {
				notifier.OnError(check, err)
			}
			notifier.AfterCheck(check)
			return checkInfo, err
		})
	}
}

type notifierLogger struct {
	Logger boshlog.Logger
}

func (n notifierLogger) BeforeCheck(check checks.Check) {
	n.Logger.Info(reflect.TypeOf(check).Name(), "Before Check ran")
}

func (n notifierLogger) AfterCheck(check checks.Check) {
	n.Logger.Info(reflect.TypeOf(check).Name(), "After Check ran")
}

func (n notifierLogger) OnError(check checks.Check, err error) {
	n.Logger.Info(reflect.TypeOf(check).Name(), "Error occurred", err)
}
