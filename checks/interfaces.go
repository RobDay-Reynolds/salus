package checks

type Check interface {
	Run() (CheckInfo, error)
}

type CheckFunc func() (CheckInfo, error)

func (cf CheckFunc) Run() (CheckInfo, error) {
	return cf()
}
