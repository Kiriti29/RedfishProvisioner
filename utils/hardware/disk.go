package redfish

import (
  "fmt"
  "github.com/redfishProvisioner/kubernetes/configmaps"
  // "github.com/redfishProvisioner/utils/hardware"
)

func main(){
    cm := configmaps.New("metal3")
    fmt.Println(cm.Get("mtn52r07c003-config"))
}
