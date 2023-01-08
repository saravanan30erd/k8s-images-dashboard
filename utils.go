package main

import (
	"os"
	"strings"
  "context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
  "k8s.io/api/apps/v1"
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

func getDeploymentMapping(deployment v1.Deployment) Deployment {
	var ready bool
	if deployment.Status.Replicas == deployment.Status.ReadyReplicas {
		ready = true
	} else {
		ready = false
	}
  return Deployment{
    Name: deployment.Name,
    Namespace: deployment.Namespace,
    Image: deployment.Spec.Template.Spec.Containers[0].Image,
		IsReady: ready,
  }
}

func getDeployments(namespace string) []Deployment {
  deploymentList, err := client.AppsV1().Deployments(namespace).List(context.TODO(),metav1.ListOptions{})
  if err != nil {
    panic(err)
  }
  deployments := make([]Deployment, len(deploymentList.Items))
  for i := 0; i < len(deployments); i++ {
		deployments[i] = getDeploymentMapping(deploymentList.Items[i])
	}
  log.Info(deployments)
	return deployments
}

func getNamespaces() []string {
	nsList, err := client.CoreV1().Namespaces().List(context.TODO(),metav1.ListOptions{})
	if err != nil {
    panic(err)
  }
	excludeNs := "kube-system, kube-public, kube-node-lease,"
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
  for i := 0; i < len(ns); i++ {
    nsDeployments[i] = NamespaceDeployments{
      Namespace: ns[i],
      Deployments: getDeployments(ns[i]),
    }
  }
  return nsDeployments
}
