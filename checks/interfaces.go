package checks

type Check interface {
	Run() error
}
