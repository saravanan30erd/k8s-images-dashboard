package main

import (
	"os"
	"strings"
  "context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NamespaceDeployments struct {
	Namespace string
	Deployments []Deployment
}

type Deployment struct {
  Name string
  Namespace string
  Image string
  IsReady bool
}

func getNamespaces() []string {
	nsList, err := client.CoreV1().Namespaces().List(context.TODO(),metav1.ListOptions{})
	if err != nil {
    panic(err)
  }
	excludeNs := "kube-system, kube-public, kube-node-lease"
	excludeNsEnv, exists := os.LookupEnv("EXCLUDE_NAMESPACES"); if exists {
		excludeNs = excludeNsEnv + ","
	}

	ns := make([]string, 0)
	for i := 0; i < len(nsList.Items); i++ {
		nsName := nsList.Items[i].Name
		if !strings.Contains(excludeNs, nsName + ",") {
			ns = append(ns, nsName)
		}
	}
  log.Info(ns)
	return ns
}

func GetNamespaceDeployments() []NamespaceDeployments {
  ns := getNamespaces()
  nsDeployments := make([]NamespaceDeployments, len(ns))
  
}
