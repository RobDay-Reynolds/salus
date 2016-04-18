package monit

import (
	"io/ioutil"
	"regexp"
	"strconv"
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
	FailedSocket FailedSocket
	DependsOn    string
}

type FileCheck struct {
	Name         string
	Path         string
	IfChanged    string
	FailedSocket FailedSocket
	Group        string
	DependsOn    string
}

type FailedSocket struct {
	SocketFile string
	Timeout    int
	NumCycles  int
	Action     string
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
			check := createProcessCheck(lines, i)
			checks = append(checks, check)
		} else if fileMatch {
			check := createFileCheck(lines, i)
			checks = append(checks, check)
		}

		i++
	}

	monitFile := MonitFile{checks}

	return monitFile, nil
}

func createProcessCheck(lines []string, startingIndex int) ProcessCheck {
	name := captureWithRegex(lines, `check process ([\w"\.]+)`, true)
	pidfile := captureWithRegex(lines, `with pidfile ([\w"/\.]+)`, true)
	startProgram := captureWithRegex(lines, `start program (.*)$`, true)
	stopProgram := captureWithRegex(lines, `stop program (.*)$`, true)
	group := captureWithRegex(lines, `group (\w+)`, true)
	dependsOn := captureWithRegex(lines, `depends on (\w+)`, true)
	failedSocket := parseFailedUnixSocket(lines)

	check := ProcessCheck{
		Name:         name,
		Pidfile:      pidfile,
		StartProgram: startProgram,
		StopProgram:  stopProgram,
		FailedSocket: failedSocket,
		Group:        group,
		DependsOn:    dependsOn,
	}

	return check
}

func createFileCheck(lines []string, startingIndex int) FileCheck {
	name := captureWithRegex(lines, `check file ([\w"\.]+)`, true)
	path := captureWithRegex(lines, `with path ([\w"/\.]+)`, true)
	ifChanged := captureWithRegex(lines, `if changed (.*)$`, true)
	group := captureWithRegex(lines, `group (\w+)`, true)
	dependsOn := captureWithRegex(lines, `depends on (\w+)`, true)
	failedSocket := parseFailedUnixSocket(lines)

	check := FileCheck{
		Name:         name,
		Path:         path,
		IfChanged:    ifChanged,
		FailedSocket: failedSocket,
		Group:        group,
		DependsOn:    dependsOn,
	}

	return check
}

func parseFailedUnixSocket(lines []string) FailedSocket {
	var startingIndex, endingIndex int
	var socketFile, timeout, numCycles, action string
	var newLines []string

	for i, line := range lines {
		newProcessCheck, err := regexp.Match("check process", []byte(line))
		if err != nil {
			// Do something
		}

		newFileCheck, err := regexp.Match("check file", []byte(line))
		if err != nil {
			// Do something
		}

		if newProcessCheck || newFileCheck {
			break
		}

		socketMatch, err := regexp.Match("if failed unixsocket", []byte(line))

		if err != nil {
			// Do something
		}

		if socketMatch {
			startingIndex = i

			newLines = append([]string{}, lines[i:]...)
			socketFile = captureWithRegex(newLines, `if failed unixsocket ([\"/a-z\.]+)`, false)
			timeout = captureWithRegex(newLines, `with timeout ([0-9]+) seconds`, false)
			numCycles = captureWithRegex(newLines, `for ([0-9]+) cycles`, false)
			action = captureWithRegex(newLines, `then ([a-z]+)`, false)

			for j, newLine := range newLines {
				thenMatch, err := regexp.Match("then ", []byte(newLine))

				if err != nil {
					// Do something
				}

				if thenMatch {
					endingIndex = i + j
				}
			}
		}
	}

	timeoutInt, err := strconv.Atoi(timeout)
	if err != nil {
		// Do something
	}

	numCyclesInt, err := strconv.Atoi(numCycles)
	if err != nil {
		// Do something
	}

	if endingIndex != 0 {
		removeElementsFromSlice(lines, startingIndex, endingIndex)
	}

	return FailedSocket{
		SocketFile: socketFile,
		Timeout:    timeoutInt,
		NumCycles:  numCyclesInt,
		Action:     action,
	}
}

func captureWithRegex(lines []string, reg string, removeLine bool) string {
	var myString string

	for i, line := range lines {
		regex := regexp.MustCompile(reg)
		values := regex.FindStringSubmatch(line)

		newProcessCheck, err := regexp.Match("check process", []byte(line))
		if err != nil {
			// Do something
		}

		newFileCheck, err := regexp.Match("check file", []byte(line))
		if err != nil {
			// Do something
		}

		if len(values) > 1 {
			myString = strings.TrimSpace(values[1])

			if removeLine {
				lines = removeElementsFromSlice(lines, i, len(lines)-1)
			}

			break
		} else if newProcessCheck || newFileCheck {
			break
		}
	}

	stripReg := regexp.MustCompile(`"([^"]*)"`)
	return stripReg.ReplaceAllString(myString, "${1}")
}

func removeElementsFromSlice(slice []string, startingIndex int, endingIndex int) []string {
	return append(slice[:startingIndex], slice[endingIndex:]...)
}
