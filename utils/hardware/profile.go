package hardware

import (
	"fmt"
	"gopkg.in/yaml.v2"
)

const (
	// DefaultProfileName is the default hardware profile to use when
	// no other profile matches.
	DefaultProfileName string = "unknown"
)

type Disks struct {
		Name string `yaml:"name"`
		RaidType string `yaml:"raid-type"`
		Disk []string `yaml:"disk"`
}

// HardwareProfile sets the hardware raid levels and nics requested by user
// on the host.
type HardwareProfile struct {
	Name	string `yaml:"name"`
	Vendor	string `yaml:"vendor"`
	Model	string `yaml:"model"`
	HardwareVersion	string `yaml:"hardware-version"`
	Bios	string `yaml:"bios"`
	Disk	[]Disks `yaml:"disks"`
	NIC []map[string]string `yaml:"network"`
}

func GetHardwareProfile(hp string) (hwp []HardwareProfile) {
	var h []HardwareProfile
  _ = yaml.Unmarshal([]byte(hp), &h)
	return h
}
