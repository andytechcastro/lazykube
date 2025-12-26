package k8s

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"lazykube/internal/domain"
	"lazykube/internal/usecase/port"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"sigs.k8s.io/yaml"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type podGateway struct {
	client  kubernetes.Interface
	config  *rest.Config
	context string
}

// NewPodGateway get the struct with the kubernetes client
func NewPodGateway(client kubernetes.Interface, config *rest.Config, cluster string) port.PodResourceGateway {
	return &podGateway{
		client:  client,
		config:  config,
		context: cluster,
	}
}

func (pg *podGateway) GetLogs(ctx context.Context, podName, namespace, containerName string) (io.ReadCloser, error) {
	logGateway := NewLogGateway(pg.client)
	return logGateway.GetLogs(ctx, podName, namespace, containerName)
}

func (pg *podGateway) Exec(ctx context.Context, podName, namespace, command, containerName string, dryRun bool, options remotecommand.StreamOptions) error {
	clientset, ok := pg.client.(*kubernetes.Clientset)
	if !ok {
		return fmt.Errorf("failed to cast client to *kubernetes.Clientset")
	}

	execGateway := NewExecGateway(clientset, pg.config)
	return execGateway.Execute(ctx, podName, containerName, namespace, command, dryRun, options)
}

func (pg *podGateway) GetAll(ctx context.Context, namespace string) ([]domain.Pod, error) {
	podList, err := pg.client.CoreV1().Pods(namespace).
		List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods in namespace %s: %w", namespace, err)
	}
	podResources := []domain.Pod{}
	for _, pod := range podList.Items {
		podResource := pg.addPodtoEntity(pod)
		podResources = append(podResources, podResource)
	}
	return podResources, nil
}

func (pg *podGateway) GetByLabels(ctx context.Context, namespace string, labelSelector map[string]string) ([]domain.Pod, error) {
	options := metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(labelSelector).String(),
	}
	podList, err := pg.client.CoreV1().Pods(namespace).
		List(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("failed to list pods with labels %v in namespace %s: %w", labelSelector, namespace, err)
	}
	podResources := []domain.Pod{}
	for _, pod := range podList.Items {
		podResource := pg.addPodtoEntity(pod)
		podResources = append(podResources, podResource)

	}
	return podResources, nil
}

func (pg *podGateway) GetYaml(ctx context.Context, namespace string, name string) ([]byte, error) {
	pod, err := pg.client.CoreV1().Pods(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get pod %s in namespace %s: %w", name, namespace, err)
	}
	pod.ManagedFields = nil
	return yaml.Marshal(pod)
}

func (pg *podGateway) GetByName(ctx context.Context, namespace string, name string) (*domain.Pod, error) {
	pod, err := pg.client.CoreV1().Pods("").Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get pod %s: %w", name, err)
	}
	podResource := pg.addPodtoEntity(*pod)
	return &podResource, nil
}

func (pg *podGateway) PortForward(namespace, podName string, ports []string, stopChan <-chan struct{}, readyChan chan struct{}) (*bytes.Buffer, *bytes.Buffer, error) {
	portForwardGateway := NewPortForwardGateway(pg.config)
	return portForwardGateway.Forward(namespace, podName, ports, stopChan, readyChan)
}

func (pg *podGateway) addPodtoEntity(pod v1.Pod) domain.Pod {
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
