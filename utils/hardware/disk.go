package main

import (
  "fmt"
  configmaps "github.com/Kiriti29/RedfishProvisioner/kubernetes"
  // "github.com/redfishProvisioner/utils/hardware"
)

func main(){
    cm := configmaps.New("metalkube")
    fmt.Println(cm.Get("mtn52r07c003-config"))
}
