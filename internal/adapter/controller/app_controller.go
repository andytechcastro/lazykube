package controller

import (
	"bytes"
	"context"
	"io"
	"lazykube/internal/domain"

	"k8s.io/client-go/tools/remotecommand"
)

// AppController Init for controller
type AppController struct {
	Pod        interface{ ControllerResource }
	Deployment interface{ ControllerResource }
	Namespace  interface{ NamespaceController }
}

type ControllerResource interface {
	GetAll(ctx context.Context, namespace string) (map[string][]map[string]string, error)
	GetAllOneContext(ctx context.Context, namespace string, context string) ([]map[string]string, error)
	GetFromManyContext(ctx context.Context, namespaces []string, contexts []string) (map[string][]map[string]string, error)
	GetYaml(ctx context.Context, namespace string, name string, context string) ([]byte, error)
	Exec(ctx context.Context, podName, namespace, context, command, containerName string, dryRun bool, options remotecommand.StreamOptions) error
	GetLogs(ctx context.Context, resourceName, namespace, context, containerName string) (io.ReadCloser, error)
	PortForward(ctx context.Context, resourceName, namespace, context string, ports []string, stopChan <-chan struct{}, readyChan chan struct{}) (*bytes.Buffer, *bytes.Buffer, error)
	GetPods(ctx context.Context, resourceName, namespace, context string) ([]domain.Pod, error)
}

type NamespaceController interface {
	GetAll(ctx context.Context, clusterContext string) ([]string, error)
}
