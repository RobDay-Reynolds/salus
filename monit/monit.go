package monit

import (
	"io/ioutil"
	"regexp"
	"strings"
)

type MonitFile struct {
	Checks []Check
}

type Check interface{}

type ProcessCheck struct {
	Pidfile      string
	StartProgram string
	StopProgram  string
	Group        string
}

func ReadMonitFile(filepath string) (MonitFile, error) {
	bytes, err := ioutil.ReadFile(filepath)

	if err != nil {
		//Do something
	}

	lines := strings.Split(string(bytes), "\n")

	pidfile := getArgForLine(lines[1], "with pidfile")
	startProgram := getArgForLine(lines[2], "start program")
	stopProgram := getArgForLine(lines[3], "stop program")
	group := getArgForLine(lines[4], "group ")

	monitFile := MonitFile{
		Checks: []Check{
			ProcessCheck{
				Pidfile:      pidfile,
				StartProgram: startProgram,
				StopProgram:  stopProgram,
				Group:        group,
			},
		},
	}

	return monitFile, nil
}

func getArgForLine(line, prefix string) string {
	myString := strings.TrimSpace(strings.Split(line, prefix)[1])

	reg := regexp.MustCompile(`"([^"]*)"`)
	return reg.ReplaceAllString(myString, "${1}")
}
