package gateways

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"lazykube/entities"
	"lazykube/usecases/gateways"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"sigs.k8s.io/yaml"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type podGateway[T entities.Pod] struct {
	client  kubernetes.Interface
	config  *rest.Config
	context string
}

// NewPodGateway get the struct with the kubernetes client
func NewPodGateway(client kubernetes.Interface, config *rest.Config, cluster string) gateways.PodResourceGateway {
	return &podGateway[entities.Pod]{
		client:  client,
		config:  config,
		context: cluster,
	}
}

func (pg *podGateway[T]) GetLogs(ctx context.Context, podName, namespace, containerName string) (io.ReadCloser, error) {
	logGateway := NewLogGateway(pg.client)
	return logGateway.GetLogs(ctx, podName, namespace, containerName)
}

func (pg *podGateway[T]) Exec(ctx context.Context, podName, namespace, command, containerName string, dryRun bool, options remotecommand.StreamOptions) error {
	clientset, ok := pg.client.(*kubernetes.Clientset)
	if !ok {
		return fmt.Errorf("failed to cast client to *kubernetes.Clientset")
	}

	execGateway := NewExecGateway(clientset, pg.config)
	return execGateway.Execute(ctx, podName, containerName, namespace, command, dryRun, options)
}

func (op *podGateway[T]) GetAll(namespace string) ([]entities.Pod, error) {
	podList, err := op.client.CoreV1().Pods(namespace).
		List(context.Background(), metav1.ListOptions{})
	if err != nil {
		fmt.Println(err)
	}
	podResources := []entities.Pod{}
	for _, pod := range podList.Items {
		podResource := op.addPodtoEntity(pod)
		podResources = append(podResources, podResource)
	}
	return podResources, nil
}

func (op *podGateway[T]) GetByLabels(namespace string, labelSelector map[string]string) ([]entities.Pod, error) {
	options := metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(labelSelector).String(),
	}
	podList, err := op.client.CoreV1().Pods(namespace).
		List(context.Background(), options)
	if err != nil {
		fmt.Println(err)
	}
	podResources := []entities.Pod{}
	for _, pod := range podList.Items {
		podResource := op.addPodtoEntity(pod)
		podResources = append(podResources, podResource)

	}
	return podResources, nil
}

func (op *podGateway[T]) GetYaml(namespace string, name string) ([]byte, error) {
	pod, err := op.client.CoreV1().Pods(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	pod.ManagedFields = nil
	return yaml.Marshal(pod)
}

func (op *podGateway[T]) GetByName(namespace string, name string) (*entities.Pod, error) {
	pod, err := op.client.CoreV1().Pods("").Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	podResource := op.addPodtoEntity(*pod)
	return &podResource, nil
}

func (pg *podGateway[T]) PortForward(namespace, podName string, ports []string, stopChan <-chan struct{}, readyChan chan struct{}) (*bytes.Buffer, *bytes.Buffer, error) {
	portForwardGateway := NewPortForwardGateway(pg.config)
	return portForwardGateway.Forward(namespace, podName, ports, stopChan, readyChan)
}

func (pg *podGateway[T]) addPodtoEntity(pod v1.Pod) entities.Pod {
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
		Context:           pg.context,
		State:             string(pod.Status.Phase),
		ContainerStatuses: statuses,
	}
	return podResource
}
