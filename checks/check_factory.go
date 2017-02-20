package checks

type CheckAdaptor func(Check) Check

func Checker(primaryCheck Check, checkAdaptors ...CheckAdaptor) Check {
	for _, adapt := range checkAdaptors {
		primaryCheck = adapt(primaryCheck)
	}
	return primaryCheck
}
