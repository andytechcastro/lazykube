package controller

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"lazykube/internal/domain"
	"lazykube/internal/usecase"
	"os"
	"text/tabwriter"

	"k8s.io/client-go/tools/remotecommand"
)

type podController struct {
	Interactor usecase.PodInteractor
}

// NewPodController return an operatorController struct
func NewPodController(interactor usecase.PodInteractor) ControllerResource {
	return &podController{
		Interactor: interactor,
	}
}

func (pC *podController) Exec(ctx context.Context, podName, namespace, context, command, containerName string, dryRun bool, options remotecommand.StreamOptions) error {
	return pC.Interactor.Exec(ctx, podName, namespace, context, command, containerName, dryRun, options)
}

func (pC *podController) GetAll(ctx context.Context, namespace string) (map[string][]map[string]string, error) {
	podLists, err := pC.Interactor.GetAll(ctx, namespace)
	if err != nil {
		return nil, err
	}

	return PodListsToMaps(podLists), nil
}

func (pC *podController) GetAllOneContext(ctx context.Context, namespace string, context string) ([]map[string]string, error) {
	pods, err := pC.Interactor.GetAllOneContext(ctx, namespace, context)
	if err != nil {
		return nil, err
	}
	return PodsToMaps(pods), nil
}

func (pC *podController) GetFromManyContext(ctx context.Context, namespaces []string, contexts []string) (map[string][]map[string]string, error) {
	pods, err := pC.Interactor.GetFromManyContext(ctx, namespaces, contexts)
	if err != nil {
		return nil, err
	}
	return PodListsToMaps(pods), nil
}

func (pD *podController) GetYaml(ctx context.Context, namespace string, name string, context string) ([]byte, error) {
	return pD.Interactor.GetYaml(ctx, namespace, name, context)
}

func (pC *podController) GetLogs(ctx context.Context, resourceName, namespace, context, containerName string) (io.ReadCloser, error) {
	return pC.Interactor.GetLogs(ctx, resourceName, namespace, context, containerName)
}

func (pC *podController) PortForward(ctx context.Context, resourceName, namespace, context string, ports []string, stopChan <-chan struct{}, readyChan chan struct{}) (*bytes.Buffer, *bytes.Buffer, error) {
	return pC.Interactor.PortForward(ctx, resourceName, namespace, context, ports, stopChan, readyChan)
}

func (pC *podController) GetPods(ctx context.Context, resourceName, namespace, context string) ([]domain.Pod, error) {
	return nil, nil
}

func (oc *podController) GetAllPodsTerminal(ctx context.Context, namespace string) error {
	podLists, err := oc.Interactor.GetAll(ctx, namespace)
	if err != nil {
		return err
	}
	tab := tabwriter.NewWriter(os.Stdout, 0, 0, 10, ' ', tabwriter.TabIndent)
	fmt.Fprintf(tab, "%s\t%s\t%s\t%s\n", "NAME", "STATUS", "CLUSTER", "NAMESPACE")
	for key, podList := range podLists {
		for _, pod := range podList {
			fmt.Fprintf(tab, "%s\t%s\t%s\t%s\n", pod.Name, pod.State, key, pod.Namespace)
		}
	}
	tab.Flush()
	return nil
}
