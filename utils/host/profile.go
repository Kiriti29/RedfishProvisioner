package host

import (
//  "fmt"
  "gopkg.in/yaml.v2"
)


type Partitions struct {
  Name string `yaml:"disk"`
  Size string `yaml:"size"`
  Bootable bool `yaml:"bootable"`
  Primary bool `yaml:"primary,omitempty"`
  FileSystem map[string]string `yaml:"filesystem"`
}

type PhysicalDevices struct {
  Disk string `yaml:"disk"`
  Partitions []Partitions `yaml:"partitions"`
}

type Storage struct {
  PhysicalDevices []PhysicalDevices `yaml:"physical-devices"`
}

type Platform struct {
  GrubConfig string `yaml:"grub_config"`
  KVMPolicy map[string]string `yaml:"kvm_policy"`
}

type HostProfile struct {
  Storage Storage `yaml:"storage"`
  Platform Platform `yaml:"platform"`
}

func GetHostProfile(hp string) (HostProfile, error) {
  var h HostProfile
  _ = yaml.Unmarshal([]byte(hp), &h)
  return h, nil
}
