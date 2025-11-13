package gateways

import (
	"context"
	"fmt"
	"io"
	"lazykube/entities"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

type LogGateway struct {
	client kubernetes.Interface
}

func NewLogGateway(client kubernetes.Interface) *LogGateway {
	return &LogGateway{client: client}
}

func (lg *LogGateway) GetLogs(ctx context.Context, podName, namespace, containerName string) (io.ReadCloser, error) {
	pod, err := lg.client.CoreV1().Pods(namespace).Get(ctx, podName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get pod %s/%s: %w", namespace, podName, err)
	}

	if containerName == "" {
		if len(pod.Spec.Containers) > 1 {
			var containerNames []string
			for _, container := range pod.Spec.Containers {
				containerNames = append(containerNames, container.Name)
			}
			return nil, &entities.ContainerSelectionError{Containers: containerNames}
		}
		containerName = pod.Spec.Containers[0].Name
	}

	podLogOpts := v1.PodLogOptions{
		Container: containerName,
		Follow:    true,
	}
	req := lg.client.CoreV1().Pods(namespace).GetLogs(podName, &podLogOpts)
	podLogs, err := req.Stream(ctx)
	if err != nil {
		return nil, fmt.Errorf("error in opening stream: %w", err)
	}

	return podLogs, nil
}
