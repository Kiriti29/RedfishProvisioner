package main

import (
  "fmt"
  configmaps "github.com/Kiriti29/RedfishProvisioner/kubernetes"
  // "github.com/redfishProvisioner/utils/hardware"
)

func main(){
    cm := configmaps.NewConfigMap("metalkube")
    label_selector := make(map[string]string)
    fmt.Println(cm.GetConfigMaps("mtn52r07c003-config", label_selector))
}
