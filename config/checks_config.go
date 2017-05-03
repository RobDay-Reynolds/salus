package config

import (
	"encoding/json"
	"github.com/monkeyherder/salus/checks"
	"github.com/monkeyherder/salus/checks/network"
	"github.com/pkg/errors"
	"time"
)

var CheckToTypeMapping map[string]func() checks.Check = map[string]func() checks.Check{
	"icmp":        func() checks.Check { return &network.IcmpCheck{} },
	"unix_socket": func() checks.Check { return &network.UnixSocketCheck{} },
	"tcp":         func() checks.Check { return &network.TcpCheck{} },
	"udp":         func() checks.Check { return &network.UdpCheck{} },
	"file":        func() checks.Check { return &checks.FileCheck{} },
	"process":     func() checks.Check { return &checks.ProcessCheck{} },
}

type CheckProperties map[string]interface{}

type CheckConfig struct {
	CheckProperties
	Type string
}

type ChecksdConfig struct {
	CheckStatusFilePath string        `json:"checkStatusPath"`
	ChecksPollTime      time.Duration `json:"checksPollTime"`
	ChecksConfig        []CheckConfig `json:"checks"`
	Checks              []checks.Check
}

func (c *ChecksdConfig) UnmarshalJSON(b []byte) error {
	type Alias ChecksdConfig

	aux := struct {
		*Alias
	}{
		Alias: (*Alias)(c),
	}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return err
	}

	for _, checkConfig := range c.ChecksConfig {
		checkForTypeFn, found := CheckToTypeMapping[checkConfig.Type]
		if !found {
			return errors.Errorf("Check config with type: '%s' is not a valid check", checkConfig.Type)
		}
		checkForType := checkForTypeFn()

		checkConfigProperties, err := json.Marshal(checkConfig.CheckProperties)
		if err != nil {
			return err
		}

		err = json.Unmarshal(checkConfigProperties, checkForType)
		if err != nil {
			return err
		}

		c.Checks = append(c.Checks, checkForType)
	}

	return nil
}
