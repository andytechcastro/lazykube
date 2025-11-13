package gateways

import (
	"context"
	"fmt"
	"lazykube/entities"
	"lazykube/usecases/gateways"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	yaml "sigs.k8s.io/yaml"
)

type deploymentGateway[T entities.Deployment] struct {
	client  kubernetes.Interface
	config  *rest.Config
	context string
}

// NewDeploymentGateway return a deploymentGateway struct
func NewDeploymentGateway(client kubernetes.Interface, config *rest.Config, cluster string) gateways.DeploymentResourceGateway {
	return &deploymentGateway[entities.Deployment]{
		client:  client,
		config:  config,
		context: cluster,
	}
}

func (dG *deploymentGateway[T]) GetPods(ctx context.Context, deploymentName, namespace string) ([]entities.Pod, error) {
	deployment, err := dG.client.AppsV1().Deployments(namespace).Get(ctx, deploymentName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment %s: %w", deploymentName, err)
	}

	selector := labels.Set(deployment.Spec.Selector.MatchLabels).String()
	podList, err := dG.client.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{LabelSelector: selector})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods for deployment %s: %w", deploymentName, err)
	}

	podResources := []entities.Pod{}
	for _, pod := range podList.Items {
		podResource := dG.addPodtoEntity(pod)
		podResources = append(podResources, podResource)
	}

	return podResources, nil
}

func (dG *deploymentGateway[T]) addPodtoEntity(pod v1.Pod) entities.Pod {
	statuses := []entities.ContainerStatuses{}
	for _, status := range pod.Status.ContainerStatuses {
		statusResource := entities.ContainerStatuses{
			Name:         status.Name,
			State:        status.State.String(),
			RestartCount: status.RestartCount,
			Image:        status.Image,
		}
		statuses = append(statuses, statusResource)
	}
	podResource := entities.Pod{
		Name:              pod.Name,
		Namespace:         pod.Namespace,
		Context:           dG.context,
		State:             string(pod.Status.Phase),
		ContainerStatuses: statuses,
	}
	return podResource
}

func (dG *deploymentGateway[T]) GetAll(namespace string) ([]entities.Deployment, error) {
	deploymentList, err := dG.client.AppsV1().Deployments(namespace).
		List(context.Background(), metav1.ListOptions{})
	if err != nil {
		fmt.Println(err)
	}
	deploymentResources := []entities.Deployment{}
	for _, deployment := range deploymentList.Items {
		deploymentResource := dG.addDeploymentEntity(deployment)
		deploymentResources = append(deploymentResources, deploymentResource)
	}
	return deploymentResources, nil
}

func (dG *deploymentGateway[T]) GetByName(namespace string, name string) (*entities.Deployment, error) {
	deployment, err := dG.client.AppsV1().Deployments(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	deploymentResource := dG.addDeploymentEntity(*deployment)

	return &deploymentResource, nil
}

func (dG *deploymentGateway[T]) GetByLabels(namespace string, label map[string]string) ([]entities.Deployment, error) {
	return []entities.Deployment{}, nil
}

func (dg *deploymentGateway[T]) addDeploymentEntity(deployment appsv1.Deployment) entities.Deployment {
	deploymentResource := entities.Deployment{
		Name:              deployment.Name,
		Namespace:         deployment.Namespace,
		Context:           dg.context,
		Replicas:          deployment.Status.Replicas,
		AvailableReplicas: deployment.Status.AvailableReplicas,
		ReadyReplicas:     deployment.Status.ReadyReplicas,
		UpdatedReplicas:   deployment.Status.UpdatedReplicas,
		MatchLabels:       deployment.Spec.Selector.MatchLabels,
	}
	return deploymentResource
}

func (dG *deploymentGateway[T]) GetYaml(namespace string, name string) ([]byte, error) {
	deployment, err := dG.client.AppsV1().Deployments(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	deployment.ManagedFields = nil
	return yaml.Marshal(deployment)
}
