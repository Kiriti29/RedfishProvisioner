package redfish

import (
  "github.com/redfishProvisioner/kubernetes/base"
   apiv1 "k8s.io/api/core/v1"
)

type SecretClient struct {
    c *apiv1.Secret
}

func New(namespace string) *SecretClient{
    clientset = base.New()
    return &SecretClient{
      c: clientset.CoreV1().Secrets("metalkube")
    }
}

func (c *SecretClient) CreateSecret(Secret *apiv1.Secret) bool {
    result, _ := c.Create(Secret)
    if result.GetObjectMeta().GetName() != ""{
        return true
    } else{
            return false
    }
}

func (j *SecretClient) DeleteSecret(name string, label_selector map[string]string) bool {
    result, _ := c.Delete(name)
    return true
}

func (j *SecretClient) GetSecrets(name string, label_selector map[string]string) bool {

}

func (j *SecretClient) GetSecretDetails(name string) bool {

}
