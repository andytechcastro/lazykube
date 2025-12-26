package k8s

import (
	"context"
	"fmt"
	"lazykube/internal/domain"
	"lazykube/internal/usecase/port"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	yaml "sigs.k8s.io/yaml"
)

type deploymentGateway struct {
	client  kubernetes.Interface
	config  *rest.Config
	context string
}

// NewDeploymentGateway return a deploymentGateway struct
func NewDeploymentGateway(client kubernetes.Interface, config *rest.Config, cluster string) port.DeploymentResourceGateway {
	return &deploymentGateway{
		client:  client,
		config:  config,
		context: cluster,
	}
}

func (pg *deploymentGateway) GetPods(ctx context.Context, deploymentName, namespace string) ([]domain.Pod, error) {
	deployment, err := pg.client.AppsV1().Deployments(namespace).Get(ctx, deploymentName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment %s in namespace %s: %w", deploymentName, namespace, err)
	}

	selector := labels.Set(deployment.Spec.Selector.MatchLabels).String()
	podList, err := pg.client.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{LabelSelector: selector})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods for deployment %s: %w", deploymentName, err)
	}

	podResources := []domain.Pod{}
	for _, pod := range podList.Items {
		podResource := pg.addPodtoEntity(pod)
		podResources = append(podResources, podResource)
	}

	return podResources, nil
}

func (pg *deploymentGateway) addPodtoEntity(pod v1.Pod) domain.Pod {
	statuses := []domain.ContainerStatuses{}
	for _, status := range pod.Status.ContainerStatuses {
		statusResource := domain.ContainerStatuses{
			Name:         status.Name,
			State:        status.State.String(),
			RestartCount: status.RestartCount,
			Image:        status.Image,
		}
		statuses = append(statuses, statusResource)
	}
	podResource := domain.Pod{
		Name:              pod.Name,
		Namespace:         pod.Namespace,
		Context:           pg.context,
		State:             string(pod.Status.Phase),
		ContainerStatuses: statuses,
	}
	return podResource
}

func (pg *deploymentGateway) GetAll(ctx context.Context, namespace string) ([]domain.Deployment, error) {
	deploymentList, err := pg.client.AppsV1().Deployments(namespace).
		List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list deployments in namespace %s: %w", namespace, err)
	}
	deploymentResources := []domain.Deployment{}
	for _, deployment := range deploymentList.Items {
		deploymentResource := pg.addDeploymentEntity(deployment)
		deploymentResources = append(deploymentResources, deploymentResource)
	}
	return deploymentResources, nil
}

func (pg *deploymentGateway) GetByName(ctx context.Context, namespace string, name string) (*domain.Deployment, error) {
	deployment, err := pg.client.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment %s in namespace %s: %w", name, namespace, err)
	}
	deploymentResource := pg.addDeploymentEntity(*deployment)

	return &deploymentResource, nil
}

func (pg *deploymentGateway) GetByLabels(ctx context.Context, namespace string, label map[string]string) ([]domain.Deployment, error) {
	return []domain.Deployment{}, nil
}

func (pg *deploymentGateway) addDeploymentEntity(deployment appsv1.Deployment) domain.Deployment {
	deploymentResource := domain.Deployment{
		Name:              deployment.Name,
		Namespace:         deployment.Namespace,
		Context:           pg.context,
		Replicas:          deployment.Status.Replicas,
		AvailableReplicas: deployment.Status.AvailableReplicas,
		ReadyReplicas:     deployment.Status.ReadyReplicas,
		UpdatedReplicas:   deployment.Status.UpdatedReplicas,
		MatchLabels:       deployment.Spec.Selector.MatchLabels,
	}
	return deploymentResource
}

func (pg *deploymentGateway) GetYaml(ctx context.Context, namespace string, name string) ([]byte, error) {
	deployment, err := pg.client.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment %s in namespace %s: %w", name, namespace, err)
	}
	deployment.ManagedFields = nil
	return yaml.Marshal(deployment)
}
