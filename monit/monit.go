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
	Name         string
	Pidfile      string
	StartProgram string
	StopProgram  string
	Group        string
}

func ReadMonitFile(filepath string) (MonitFile, error) {
	bytes, err := ioutil.ReadFile(filepath)

	if err != nil {
		// Do something
	}

	lines := strings.Split(string(bytes), "\n")

	checks := []Check{}

	i := 0
	for _, line := range lines {
		match, err := regexp.Match("check process", []byte(line))

		if err != nil {
			// Do something
		}

		if match {
			check := createProcessCheck(lines, i)
			checks = append(checks, check)

			lines = append(lines[:i], lines[i+5:]...)
		}

		i++
	}

	monitFile := MonitFile{checks}

	return monitFile, nil
}

func createProcessCheck(lines []string, startingIndex int) ProcessCheck {
	name := getArgForLine(lines, "check process")
	pidfile := getArgForLine(lines, "with pidfile")
	startProgram := getArgForLine(lines, "start program")
	stopProgram := getArgForLine(lines, "stop program")
	group := getArgForLine(lines, "group ")

	check := ProcessCheck{
		Name:         name,
		Pidfile:      pidfile,
		StartProgram: startProgram,
		StopProgram:  stopProgram,
		Group:        group,
	}

	return check
}

func getArgForLine(lines []string, prefix string) string {
	var myString string

	for _, line := range lines {
		match, err := regexp.Match(prefix, []byte(line))

		if err != nil {
			// Do something
		}

		if match {
			myString = strings.TrimSpace(strings.Split(line, prefix)[1])

			break
		}
	}

	reg := regexp.MustCompile(`"([^"]*)"`)
	return reg.ReplaceAllString(myString, "${1}")
}
