package adaptors

import "github.com/monkeyherder/moirai/checks"

type Summary struct {
	CheckInfos checks.CheckInfo
	CheckError error
}

type CheckSummary struct {
	CheckSummary []Summary
}

func PersistCheckSummary(pathToCheckSummary string) checks.CheckAdaptor {
	return func(check checks.Check) checks.Check {
		return check
		//		return checks.CheckFunc(func() (checks.CheckInfo, error) {
		//
		//			_, err := check.Run()
		//			if err != nil {
		//
		//			}
		//
		//			return checks.CheckInfo{}, err
		//		})
	}
}
