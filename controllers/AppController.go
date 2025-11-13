package controllers

import (
	"bytes"
	"context"
	"io"
	"lazykube/entities"

	"k8s.io/client-go/tools/remotecommand"
)

// AppController Init for controller
type AppController struct {
	Pod        interface{ ControllerResource }
	Deployment interface{ ControllerResource }
	Namespace  interface{ NamespaceController }
}

type ControllerResource interface {
	GetAll(namespace string) (map[string][]map[string]string, error)
	GetAllOneContext(namespace string, context string) ([]map[string]string, error)
	GetFromManyContext(namespaces []string, contexts []string) (map[string][]map[string]string, error)
	GetYaml(string, string, string) ([]byte, error)
	Exec(ctx context.Context, podName, namespace, context, command, containerName string, dryRun bool, options remotecommand.StreamOptions) error
	GetLogs(ctx context.Context, resourceName, namespace, context, containerName string) (io.ReadCloser, error)
	PortForward(ctx context.Context, resourceName, namespace, context string, ports []string, stopChan <-chan struct{}, readyChan chan struct{}) (*bytes.Buffer, *bytes.Buffer, error)
	GetPods(ctx context.Context, resourceName, namespace, context string) ([]entities.Pod, error)
}
