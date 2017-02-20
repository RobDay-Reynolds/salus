package checks

type Check interface {
	Run() error
}

type CheckFunc func() error

func (cf CheckFunc) Run() error {
	return cf()
}
