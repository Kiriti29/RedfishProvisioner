package github.com/Kiriti29/RedfishProvisioner/kubernetes

import (
  //"github.com/redfishProvisioner/kubernetes/base"
   apiv1 "k8s.io/api/core/v1"
)

type PodClient struct {
    c *apiv1.Pod
}

func New(namespace string) *PodClient{
    clientset = base.New()
    return &PodClient{
      c: clientset.CoreV1().Pods("metalkube")
    }
}

func (c *PodClient) CreatePod(Pod *apiv1.Pod) bool {
    result, _ := c.Create(Pod)
    if result.GetObjectMeta().GetName() != ""{
        return true
    } else{
            return false
    }
}

func (j *PodClient) DeletePod(name string, label_selector map[string]string) bool {
    result, _ := c.Delete(name)
    return true
}

func (j *PodClient) GetPods(name string, label_selector map[string]string) bool {

}

func (j *PodClient) GetPodDetails(name string) bool {

}
