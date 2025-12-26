package controller

import (
	"fmt"
	"lazykube/internal/domain"
)

func PodToMap(pod domain.Pod) map[string]string {
	return map[string]string{
		"name":      pod.Name,
		"namespace": pod.Namespace,
		"cluster":   pod.Context,
		"status":    pod.State,
	}
}

func PodsToMaps(pods []domain.Pod) []map[string]string {
	result := make([]map[string]string, len(pods))
	for i, pod := range pods {
		result[i] = PodToMap(pod)
	}
	return result
}

func PodListsToMaps(podLists map[string][]domain.Pod) map[string][]map[string]string {
	result := make(map[string][]map[string]string)
	for cluster, pods := range podLists {
		result[cluster] = PodsToMaps(pods)
	}
	return result
}

func DeploymentToMap(deployment domain.Deployment) map[string]string {
	return map[string]string{
		"name":      deployment.Name,
		"namespace": deployment.Namespace,
		"cluster":   deployment.Context,
		"replicas":  fmt.Sprintf("%d/%d", deployment.ReadyReplicas, deployment.Replicas),
	}
}

func DeploymentsToMaps(deployments []domain.Deployment) []map[string]string {
	result := make([]map[string]string, len(deployments))
	for i, deployment := range deployments {
		result[i] = DeploymentToMap(deployment)
	}
	return result
}

func DeploymentListsToMaps(deploymentLists map[string][]domain.Deployment) map[string][]map[string]string {
	result := make(map[string][]map[string]string)
	for cluster, deployments := range deploymentLists {
		result[cluster] = DeploymentsToMaps(deployments)
	}
	return result
}
