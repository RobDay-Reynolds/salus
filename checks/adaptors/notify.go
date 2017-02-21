package adaptors

import (
	"fmt"

	"github.com/monkeyherder/moirai/checks"
)

func NewNotifierLogger() checks.CheckAdaptor {
	return Notify(notifierLogger{})
}

type Notifier interface {
	BeforeCheck(checks.Check)
	AfterCheck(checks.Check)
	OnError(checks.Check, error)
}

func Notify(notifier Notifier) checks.CheckAdaptor {
	return func(check checks.Check) checks.Check {
		return checks.CheckFunc(func() (string, string, error) {
			notifier.BeforeCheck(check)
			_, _, err := check.Run()
			if err != nil {
				notifier.OnError(check, err)
			}
			notifier.AfterCheck(check)
			return "", "", err
		})
	}
}

type notifierLogger struct{}

func (notifierLogger) BeforeCheck(check checks.Check) {
	fmt.Println("Before Check ran")
}

func (notifierLogger) AfterCheck(check checks.Check) {
	fmt.Println("After Check ran")
}

func (notifierLogger) OnError(chgeck checks.Check, err error) {
	fmt.Println(fmt.Sprintf("Error occurred: %v", err))
}
