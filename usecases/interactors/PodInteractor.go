package interactors

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"lazykube/entities"
	"lazykube/usecases/gateways"
	"sync"

	"k8s.io/client-go/tools/remotecommand"
)

type podInteractor struct {
	PodRepo map[string]gateways.PodResourceGateway
}

// PodInteractor is the interface for connect with this struct
type PodInteractor interface {
	GetAll(string) (map[string][]entities.Pod, error)
	GetAllOneContext(string, string) ([]entities.Pod, error)
	GetFromManyContext([]string, []string) (map[string][]entities.Pod, error)
	GetYaml(string, string, string) ([]byte, error)
	Exec(ctx context.Context, podName, namespace, context, command, containerName string, dryRun bool, options remotecommand.StreamOptions) error
	GetLogs(ctx context.Context, podName, namespace, context, containerName string) (io.ReadCloser, error)
	PortForward(ctx context.Context, podName, namespace, context string, ports []string, stopChan <-chan struct{}, readyChan chan struct{}) (*bytes.Buffer, *bytes.Buffer, error)
}

// NewPodInteractor return an struct of tyoe operatorInteractor
func NewPodInteractor(podRepo map[string]gateways.PodResourceGateway) PodInteractor {
	return &podInteractor{
		PodRepo: podRepo,
	}
}

func (pi *podInteractor) GetLogs(ctx context.Context, podName, namespace, context, containerName string) (io.ReadCloser, error) {
	gateway := pi.PodRepo[context]
	if gateway == nil {
		return nil, fmt.Errorf("no gateway found for context: %s", context)
	}
	return gateway.GetLogs(ctx, podName, namespace, containerName)
}

func (pi *podInteractor) Exec(ctx context.Context, podName, namespace, context, command, containerName string, dryRun bool, options remotecommand.StreamOptions) error {
	gateway := pi.PodRepo[context]
	if gateway == nil {
		return fmt.Errorf("no gateway found for context: %s", context)
	}
	return gateway.Exec(ctx, podName, namespace, command, containerName, dryRun, options)
}

func (oi *podInteractor) GetAll(namespace string) (map[string][]entities.Pod, error) {
	podLists := map[string][]entities.Pod{}
	wg := sync.WaitGroup{}
	var err error
	for key, repo := range oi.PodRepo {
		wg.Add(1)
		go func(namespace string, key string, repo gateways.PodResourceGateway) {
			podLists[key], err = repo.GetAll(namespace)
			if err != nil {
				fmt.Println(err)
			}
			defer wg.Done()
		}(namespace, key, repo)
	}
	wg.Wait()
	return podLists, nil
}

func (oi *podInteractor) GetFromManyContext(namespaces []string, contexts []string) (map[string][]entities.Pod, error) {
	podLists := map[string][]entities.Pod{}
	wg := sync.WaitGroup{}
	for _, context := range contexts {
		for _, namespace := range namespaces {
			wg.Add(1)
			repo := oi.PodRepo[context]
			go func(namespace string, context string, repo gateways.PodResourceGateway) {
				podNamespaceLists, err := repo.GetAll(namespace)
				if err != nil {
					fmt.Println(err)
				}
				podLists[context] = append(podLists[context], podNamespaceLists...)
				defer wg.Done()
			}(namespace, context, repo)
		}
	}
	wg.Wait()
	return podLists, nil
}

func (dI *podInteractor) GetAllOneContext(namespace string, context string) ([]entities.Pod, error) {
	deploymentLists := []entities.Pod{}
	var err error
	deploymentLists, err = dI.PodRepo[context].GetAll(namespace)
	if err != nil {
		fmt.Println(err)
	}
	return deploymentLists, nil
}

func (pi *podInteractor) PortForward(ctx context.Context, podName, namespace, context string, ports []string, stopChan <-chan struct{}, readyChan chan struct{}) (*bytes.Buffer, *bytes.Buffer, error) {
	gateway := pi.PodRepo[context]
	if gateway == nil {
		return nil, nil, fmt.Errorf("no gateway found for context: %s", context)
	}
	return gateway.PortForward(namespace, podName, ports, stopChan, readyChan)
}

func (dI *podInteractor) GetYaml(namespace string, name string, context string) ([]byte, error) {
	return dI.PodRepo[context].GetYaml(namespace, name)
}
