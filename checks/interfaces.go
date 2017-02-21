package checks

type Check interface {
	Run() (string, string, error)
}

type CheckFunc func() (string, string, error)

func (cf CheckFunc) Run() (string, string, error) {
	return cf()
}
