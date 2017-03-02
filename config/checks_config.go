package config

import (
	"github.com/monkeyherder/moirai/checks/network"
	"time"
)

type ChecksdConfig struct {
	ChecksPollTime time.Duration       `json:"checksPollTime"`
	IcmpChecks     []network.IcmpCheck `json:"icmpChecks"`
}
