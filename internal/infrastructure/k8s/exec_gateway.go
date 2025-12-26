package k8s

import (
	"context"
	"fmt"
	"lazykube/internal/domain"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

// IExecGateway defines the interface for executing commands in a pod.
type IExecGateway interface {
	Execute(ctx context.Context, podName, containerName, namespace, command string, options remotecommand.StreamOptions) error
}

// ExecGateway implements the IExecGateway interface.
type ExecGateway struct {
	clientset *kubernetes.Clientset
	config    *rest.Config
}

// NewExecGateway creates a new ExecGateway.
func NewExecGateway(clientset *kubernetes.Clientset, config *rest.Config) *ExecGateway {
	return &ExecGateway{
		clientset: clientset,
		config:    config,
	}
}

// Execute establishes a SPDY connection to a pod and executes a command.
func (g *ExecGateway) Execute(ctx context.Context, podName, containerName, namespace, command string, dryRun bool, options remotecommand.StreamOptions) error {
	pod, err := g.clientset.CoreV1().Pods(namespace).Get(ctx, podName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get pod %s/%s: %w", namespace, podName, err)
	}

	if containerName == "" {
		if len(pod.Spec.Containers) > 1 {
			var containerNames []string
			for _, container := range pod.Spec.Containers {
				containerNames = append(containerNames, container.Name)
			}
			return &domain.ContainerSelectionError{Containers: containerNames}
		}
		if len(pod.Spec.Containers) > 0 {
			containerName = pod.Spec.Containers[0].Name
		}
	}

	// Dry run mode
	if dryRun {
		return nil
	}

	req := g.clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		Param("container", containerName).
		Param("command", command).
		Param("stdin", "true").
		Param("stdout", "true").
		Param("stderr", "true").
		Param("tty", "true").
		Param("env", "TERM=xterm-256color")

	executor, err := remotecommand.NewSPDYExecutor(g.config, "POST", req.URL())
	if err != nil {
		return err
	}

	return executor.StreamWithContext(ctx, options)
}
