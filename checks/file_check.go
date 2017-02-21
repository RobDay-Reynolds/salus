package checks

import "os"

type FileCheck struct {
	Name      string
	Path      string
	IfChanged string
	Group     string
	DependsOn string
}

func (f FileCheck) Run() (string, string, error) {
	_, err := os.Stat(f.Path)

	return "", "", err
}
