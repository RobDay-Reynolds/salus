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
	FailedSocket FailedSocket
	FailedHost   FailedHost
	Group        string
	DependsOn    string
}

type FileCheck struct {
	Name         string
	Path         string
	IfChanged    string
	FailedSocket FailedSocket
	FailedHost   FailedHost
	Group        string
	DependsOn    string
}

type FailedSocket struct {
	SocketFile string
	Timeout    int
	NumCycles  int
	Action     string
}

type FailedHost struct {
	Host      string
	Port      string
	Protocol  string
	Timeout   int
	NumCycles int
	Action    string
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
	failedHost := parseFailedHost(lines)

	check := ProcessCheck{
		Name:         name,
		Pidfile:      pidfile,
		StartProgram: startProgram,
		StopProgram:  stopProgram,
		FailedSocket: failedSocket,
		FailedHost:   failedHost,
		Group:        group,
		DependsOn:    dependsOn,
	}

	return check
}

func createFileCheck(lines []string, startingIndex int) FileCheck {
	name := captureWithRegex(lines, `check file ([\w"\.]+)`, true)
	failedHost := parseFailedHost(lines)
	failedSocket := parseFailedUnixSocket(lines)

	path := captureWithRegex(lines, `with path ([\w"/\.]+)`, true)
	ifChanged := captureWithRegex(lines, `if changed (.*)$`, true)
	group := captureWithRegex(lines, `group (\w+)`, true)
	dependsOn := captureWithRegex(lines, `depends on (\w+)`, true)

	check := FileCheck{
		Name:         name,
		Path:         path,
		IfChanged:    ifChanged,
		FailedSocket: failedSocket,
		FailedHost:   failedHost,
		Group:        group,
		DependsOn:    dependsOn,
	}

	return check
}

func parseFailedUnixSocket(lines []string) FailedSocket {
	values := parseGroupBlock(
		lines,
		"socketFile",
		map[string]string{
			"socketFile": `if failed unixsocket (["/\w\.]+)`,
			"timeout":    `with timeout ([0-9]+) seconds`,
			"numCycles":  `for ([0-9]+) cycles`,
			"action":     `then ([a-z]+)`,
		},
	)

	socketFile := values["socketFile"]
	timeout := values["timeout"]
	numCycles := values["numCycles"]
	action := values["action"]

	timeoutInt, err := strconv.Atoi(timeout)
	if err != nil {
		// Do something
	}

	numCyclesInt, err := strconv.Atoi(numCycles)
	if err != nil {
		// Do something
	}

	return FailedSocket{
		SocketFile: socketFile,
		Timeout:    timeoutInt,
		NumCycles:  numCyclesInt,
		Action:     action,
	}
}

func parseFailedHost(lines []string) FailedHost {
	values := parseGroupBlock(
		lines,
		"host",
		map[string]string{
			"host":      `if failed host ([\w\.]+)`,
			"port":      `port ([\d]+)`,
			"protocol":  `protocol ([\w]+)`,
			"timeout":   `with timeout ([0-9]+) seconds`,
			"numCycles": `for ([0-9]+) cycles`,
			"action":    `then ([a-z]+)`,
		},
	)

	host := values["host"]
	port := values["port"]
	protocol := values["protocol"]
	timeout := values["timeout"]
	numCycles := values["numCycles"]
	action := values["action"]

	timeoutInt, err := strconv.Atoi(timeout)
	if err != nil {
		// Do something
	}

	numCyclesInt, err := strconv.Atoi(numCycles)
	if err != nil {
		// Do something
	}

	return FailedHost{
		Host:      host,
		Port:      port,
		Protocol:  protocol,
		Timeout:   timeoutInt,
		NumCycles: numCyclesInt,
		Action:    action,
	}
}

func parseGroupBlock(lines []string, keyRegex string, regexes map[string]string) map[string]string {
	var startingIndex, endingIndex int
	var newLines []string
	values := map[string]string{}

	startingRegex := regexp.MustCompile(regexes[keyRegex])

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

		match := startingRegex.Match([]byte(line))

		if match {
			startingIndex = i

			newLines = append([]string{}, lines[i:]...)

			for key, regex := range regexes {
				values[key] = captureWithRegex(newLines, regex, false)
			}

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

	if endingIndex != 0 {
		removeElementsFromSlice(lines, startingIndex, endingIndex)
	}

	return values
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
