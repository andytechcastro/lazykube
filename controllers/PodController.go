package controllers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"lazykube/entities"
	"lazykube/usecases/interactors"
	"os"
	"text/tabwriter"

	"log"

	"k8s.io/client-go/tools/remotecommand"
)

type podController struct {
	Interactor interactors.PodInteractor
}

// NewPodController return an operatorController struct
func NewPodController(interactor interactors.PodInteractor) ControllerResource {
	return &podController{
		Interactor: interactor,
	}
}

func (pC *podController) Exec(ctx context.Context, podName, namespace, context, command, containerName string, dryRun bool, options remotecommand.StreamOptions) error {
	return pC.Interactor.Exec(ctx, podName, namespace, context, command, containerName, dryRun, options)
}

func (pC *podController) GetAll(namespace string) (map[string][]map[string]string, error) {
	podLists, err := pC.Interactor.GetAll(namespace)
	if err != nil {
		fmt.Println(err)
		log.Panic()
	}

	return ResourceToData[map[string][]map[string]string](podLists)
}

func (pC *podController) GetAllOneContext(namespace string, context string) ([]map[string]string, error) {
	pods, err := pC.Interactor.GetAllOneContext(namespace, context)
	if err != nil {
		return nil, err
	}
	return ResourceToData[[]map[string]string](pods)
}

func (pC *podController) GetFromManyContext(namespaces []string, contexts []string) (map[string][]map[string]string, error) {
	pods, err := pC.Interactor.GetFromManyContext(namespaces, contexts)
	if err != nil {
		return nil, err
	}
	return ResourceToData[map[string][]map[string]string](pods)
}

func (pD *podController) GetYaml(namespace string, name string, context string) ([]byte, error) {
	return pD.Interactor.GetYaml(namespace, name, context)
}

func (pC *podController) GetLogs(ctx context.Context, resourceName, namespace, context, containerName string) (io.ReadCloser, error) {
	return pC.Interactor.GetLogs(ctx, resourceName, namespace, context, containerName)
}

func (pC *podController) PortForward(ctx context.Context, resourceName, namespace, context string, ports []string, stopChan <-chan struct{}, readyChan chan struct{}) (*bytes.Buffer, *bytes.Buffer, error) {
	return pC.Interactor.PortForward(ctx, resourceName, namespace, context, ports, stopChan, readyChan)
}

func (pC *podController) GetPods(ctx context.Context, resourceName, namespace, context string) ([]entities.Pod, error) {
	return nil, nil
}

func (oc *podController) GetAllPodsTerminal(namespace string) {
	podLists, err := oc.Interactor.GetAll(namespace)
	if err != nil {
		fmt.Println(err)
		log.Panic()
	}
	tab := tabwriter.NewWriter(os.Stdout, 0, 0, 10, ' ', tabwriter.TabIndent)
	fmt.Fprintf(tab, "%s\t%s\t%s\t%s\n", "NAME", "STATUS", "CLUSTER", "NAMESPACE")
	for key, podList := range podLists {
		for _, pod := range podList {
			fmt.Fprintf(tab, "%s\t%s\t%s\t%s\n", pod.Name, pod.State, key, pod.Namespace)
		}
	}
	tab.Flush()
}
