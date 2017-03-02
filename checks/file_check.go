package checks

import "os"

type FileCheck struct {
	Name      string
	Path      string
	IfChanged string
	Group     string
	DependsOn string
}

func (f FileCheck) Run() (CheckInfo, error) {
	_, err := os.Stat(f.Path)

	return CheckInfo{}, err
}
