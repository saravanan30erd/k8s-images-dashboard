package main

import (
	"context"
	"k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"strings"
)

func concatMultipleSlices[T any](slices [][]T) []T {
	var totalLen int
	for _, s := range slices {
		totalLen += len(s)
	}
	result := make([]T, totalLen)
	var i int
	for _, s := range slices {
		i += copy(result[i:], s)
	}
	return result
}

type NamespaceWorkloads struct {
	Namespace string
	Workloads []Workload
}

type Workload struct {
	Name    string
	Type    string
	Image   string
	IsReady bool
}

func getDeploymentMapping(deployment v1.Deployment) Workload {
	var ready bool
	if deployment.Status.Replicas == deployment.Status.ReadyReplicas {
		ready = true
	} else {
		ready = false
	}
	return Workload{
		Name:    deployment.Name,
		Type:    "Deployment",
		Image:   deployment.Spec.Template.Spec.Containers[0].Image,
		IsReady: ready,
	}
}

func getDaemonSetMapping(daemonset v1.DaemonSet) Workload {
	var ready bool
	if daemonset.Status.DesiredNumberScheduled == daemonset.Status.NumberReady {
		ready = true
	} else {
		ready = false
	}
	return Workload{
		Name:    daemonset.Name,
		Type:    "DaemonSet",
		Image:   daemonset.Spec.Template.Spec.Containers[0].Image,
		IsReady: ready,
	}
}

func getStatefulSetMapping(statefulSet v1.StatefulSet) Workload {
	var ready bool
	if statefulSet.Status.Replicas == statefulSet.Status.ReadyReplicas {
		ready = true
	} else {
		ready = false
	}
	return Workload{
		Name:    statefulSet.Name,
		Type:    "StatefulSet",
		Image:   statefulSet.Spec.Template.Spec.Containers[0].Image,
		IsReady: ready,
	}
}

func getWorkloads(namespace string) []Workload {
	var workloads []Workload
	deploymentList, err := client.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	deployments := make([]Workload, len(deploymentList.Items))
	for i := 0; i < len(deployments); i++ {
		deployments[i] = getDeploymentMapping(deploymentList.Items[i])
	}

	onlyDeployment, _ := os.LookupEnv("SHOW_ONLY_DEPLOYMENT")
	if onlyDeployment == "true" {
		workloads = deployments
	} else {
		statefulSetsList, err := client.AppsV1().StatefulSets(namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err)
		}
		statefulSets := make([]Workload, len(statefulSetsList.Items))
		for i := 0; i < len(statefulSets); i++ {
			statefulSets[i] = getStatefulSetMapping(statefulSetsList.Items[i])
		}

		daemonSetsList, err := client.AppsV1().DaemonSets(namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err)
		}
		daemonSets := make([]Workload, len(daemonSetsList.Items))
		for i := 0; i < len(daemonSets); i++ {
			daemonSets[i] = getDaemonSetMapping(daemonSetsList.Items[i])
		}

		workloads = concatMultipleSlices([][]Workload{deployments, daemonSets, statefulSets})
	}
	return workloads
}

func getNamespaces() []string {
	nsList, err := client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	excludeNs := "kube-system, kube-public, kube-node-lease, default,"
	excludeNsEnv, exists := os.LookupEnv("EXCLUDE_NAMESPACES")
	if exists {
		excludeNs = excludeNsEnv + ","
	}

	ns := make([]string, 0)
	for i := 0; i < len(nsList.Items); i++ {
		nsName := nsList.Items[i].Name
		if !strings.Contains(excludeNs, nsName+",") {
			ns = append(ns, nsName)
		}
	}
	return ns
}

func GetNamespaceWorkloads() []NamespaceWorkloads {
	ns := getNamespaces()
	nsWorkloads := make([]NamespaceWorkloads, len(ns))
	for i := 0; i < len(ns); i++ {
		nsWorkloads[i] = NamespaceWorkloads{
			Namespace: ns[i],
			Workloads: getWorkloads(ns[i]),
		}
	}
	//log.Info(nsWorkloads)
	return nsWorkloads
}
