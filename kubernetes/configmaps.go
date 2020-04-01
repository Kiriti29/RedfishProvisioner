package github.com/Kiriti29/RedfishProvisioner/kubernetes

import (
   // base "github.com/RedfishProvisioner/kubernetes/base"
   apiv1 "k8s.io/api/core/v1"
)

type ConfigMapClient struct {
    c *apiv1.ConfigMap
}

func New(namespace string) *ConfigMapClient{
    clientset = base.New()
    return &ConfigMapClient{
      c: clientset.CoreV1().ConfigMaps("metalkube")
    }
}

func (c *ConfigMapClient) CreateConfigMap(ConfigMap *apiv1.ConfigMap) bool {
    result, _ := c.Create(ConfigMap)
    if result.GetObjectMeta().GetName() != ""{
        return true
    } else{
            return false
    }
}

func (j *ConfigMapClient) DeleteConfigMap(name string, label_selector map[string]string) bool {
    result, _ := c.Delete(name)
    return true
}

func (j *ConfigMapClient) GetConfigMaps(name string, label_selector map[string]string) bool {
  result, _ := c.Get(name)
  return result
}

func (j *ConfigMapClient) GetConfigMapDetails(name string) bool {

}
