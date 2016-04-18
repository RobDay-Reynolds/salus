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
	DependsOn    string
}

type FileCheck struct {
	Name      string
	Path      string
	IfChanged string
	Group     string
	DependsOn string
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
		processMatch, err := regexp.Match("check process", []byte(line))
		fileMatch, err := regexp.Match("check file", []byte(line))

		if err != nil {
			// Do something
		}

		if processMatch {
			check, numLines := createProcessCheck(lines, i)
			checks = append(checks, check)

			lines = append([]string{}, lines[numLines:]...)
		} else if fileMatch {
			check, numLines := createFileCheck(lines, i)
			checks = append(checks, check)

			lines = append([]string{}, lines[numLines:]...)
		}

		i++
	}

	monitFile := MonitFile{checks}

	return monitFile, nil
}

func createProcessCheck(lines []string, startingIndex int) (ProcessCheck, int) {
	name := getArgForLine(lines, "check process")
	pidfile := getArgForLine(lines, "with pidfile")
	startProgram := getArgForLine(lines, "start program")
	stopProgram := getArgForLine(lines, "stop program")
	group := getArgForLine(lines, "group ")
	dependsOn := getArgForLine(lines, "depends on ")

	values := []string{
		name,
		pidfile,
		startProgram,
		stopProgram,
		group,
		dependsOn,
	}
	numMatches := 0

	for _, value := range values {
		if value != "" {
			numMatches++
		}
	}

	check := ProcessCheck{
		Name:         name,
		Pidfile:      pidfile,
		StartProgram: startProgram,
		StopProgram:  stopProgram,
		Group:        group,
		DependsOn:    dependsOn,
	}

	return check, numMatches
}

func createFileCheck(lines []string, startingIndex int) (FileCheck, int) {
	name := getArgForLine(lines, "check file")
	path := getArgForLine(lines, "with path")
	ifChanged := getArgForLine(lines, "if changed")
	group := getArgForLine(lines, "group ")
	dependsOn := getArgForLine(lines, "depends on ")

	values := []string{
		name,
		path,
		ifChanged,
		group,
		dependsOn,
	}
	numMatches := 0

	for _, value := range values {
		if value != "" {
			numMatches++
		}
	}

	check := FileCheck{
		Name:      name,
		Path:      path,
		IfChanged: ifChanged,
		Group:     group,
		DependsOn: dependsOn,
	}

	return check, numMatches
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
