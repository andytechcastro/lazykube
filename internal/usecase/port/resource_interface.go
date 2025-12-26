package port

import (
	"bytes"
	"context"
	"io"
	"lazykube/internal/domain"

	"k8s.io/client-go/tools/remotecommand"
)

// ResourceGateway is a generic interface for accessing Kubernetes resources.
type ResourceGateway[T Resource] interface {
	GetAll(ctx context.Context, namespace string) ([]T, error)
	GetByName(ctx context.Context, namespace string, name string) (*T, error)
	GetByLabels(ctx context.Context, namespace string, label map[string]string) ([]T, error)
	GetYaml(ctx context.Context, namespace string, name string) ([]byte, error)
}

// PodResourceGateway defines operations specific to Pods, including streaming.
type PodResourceGateway interface {
	ResourceGateway[domain.Pod]
	Exec(ctx context.Context, podName, namespace, command, containerName string, dryRun bool, options remotecommand.StreamOptions) error
	GetLogs(ctx context.Context, podName, namespace, containerName string) (io.ReadCloser, error)
	PortForward(namespace, podName string, ports []string, stopChan <-chan struct{}, readyChan chan struct{}) (*bytes.Buffer, *bytes.Buffer, error)
}

// DeploymentResourceGateway defines operations specific to Deployments.
type DeploymentResourceGateway interface {
	ResourceGateway[domain.Deployment]
	GetPods(ctx context.Context, deploymentName, namespace string) ([]domain.Pod, error)
}

type Resource interface {
	domain.Deployment | domain.Pod | string
}
