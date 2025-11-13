package gateways

import (
	"bytes"
	"context"
	"io"
	"lazykube/entities"

	"k8s.io/client-go/tools/remotecommand"
)

// ResourceGateway is a generic interface for accessing Kubernetes resources.
type ResourceGateway[T Resource] interface {
	GetAll(namespace string) ([]T, error)
	GetByName(namespace string, name string) (*T, error)
	GetByLabels(namespace string, label map[string]string) ([]T, error)
	GetYaml(namespace string, name string) ([]byte, error)
}

// PodResourceGateway defines operations specific to Pods, including streaming.
type PodResourceGateway interface {
	ResourceGateway[entities.Pod]
	Exec(ctx context.Context, podName, namespace, command, containerName string, dryRun bool, options remotecommand.StreamOptions) error
	GetLogs(ctx context.Context, podName, namespace, containerName string) (io.ReadCloser, error)
	PortForward(namespace, podName string, ports []string, stopChan <-chan struct{}, readyChan chan struct{}) (*bytes.Buffer, *bytes.Buffer, error)
}

// DeploymentResourceGateway defines operations specific to Deployments.
type DeploymentResourceGateway interface {
	ResourceGateway[entities.Deployment]
	GetPods(ctx context.Context, deploymentName, namespace string) ([]entities.Pod, error)
}

type Resource interface {
	entities.Deployment | entities.Pod | string
}
