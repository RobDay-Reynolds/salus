package config

import (
	"encoding/json"
	"github.com/monkeyherder/moirai/checks"
	"github.com/monkeyherder/moirai/checks/network"
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
	ChecksPollTime time.Duration `json:"checksPollTime"`
	ChecksConfig   []CheckConfig `json:"checks"`
	Checks         []checks.Check
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
		check := CheckToTypeMapping[checkConfig.Type]()

		checkConfigProperties, err := json.Marshal(checkConfig.CheckProperties)
		if err != nil {
			return err
		}

		err = json.Unmarshal(checkConfigProperties, check)
		if err != nil {
			return err
		}

		c.Checks = append(c.Checks, check)
	}

	return nil
}
