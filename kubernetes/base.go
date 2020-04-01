package github.com/Kiriti29/RedfishProvisioner/kubernetes

import (
  "k8s.io/client-go/kubernetes"
  "k8s.io/client-go/rest"
)
type Kubernetes struct {
    kube *kubernetes.Clientset
}

(p *Kubernetes) New() *kubernetes.Clientset{
    // creates the in-cluster config
    config, err := rest.InClusterConfig()
    if err != nil {
      panic(err.Error())
    }
    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
      panic(err.Error())
    }
    return &Kubernetes{
      kube: clientset
    }
}
