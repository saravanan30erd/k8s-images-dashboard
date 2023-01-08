package main

import (
    "flag"
    "path/filepath"
    "k8s.io/client-go/util/homedir"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
)

var client *kubernetes.Clientset

func init() {
  var kubeconfig *string
  if home := homedir.HomeDir(); home != "" {
      kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "kubeconfig file")
  } else {
      kubeconfig = flag.String("kubeconfig", "", "kubeconfig file")
  }
  flag.Parse()
  config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
  if err != nil {
      panic(err)
  }
  clientset, err := kubernetes.NewForConfig(config)
  if err != nil {
      panic(err)
  }
  client = clientset
}
