package usecase

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"lazykube/internal/domain"
	"lazykube/internal/usecase/port"
	"sync"

	"k8s.io/client-go/tools/remotecommand"
)

type podInteractor struct {
	PodRepo map[string]port.PodResourceGateway
}

// PodInteractor is the interface for connect with this struct
type PodInteractor interface {
	GetAll(context.Context, string) (map[string][]domain.Pod, error)
	GetAllOneContext(context.Context, string, string) ([]domain.Pod, error)
	GetFromManyContext(context.Context, []string, []string) (map[string][]domain.Pod, error)
	GetYaml(context.Context, string, string, string) ([]byte, error)
	Exec(ctx context.Context, podName, namespace, context, command, containerName string, dryRun bool, options remotecommand.StreamOptions) error
	GetLogs(ctx context.Context, podName, namespace, context, containerName string) (io.ReadCloser, error)
	PortForward(ctx context.Context, podName, namespace, context string, ports []string, stopChan <-chan struct{}, readyChan chan struct{}) (*bytes.Buffer, *bytes.Buffer, error)
}

// NewPodInteractor return an struct of tyoe operatorInteractor
func NewPodInteractor(podRepo map[string]port.PodResourceGateway) PodInteractor {
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

func (pi *podInteractor) GetAll(ctx context.Context, namespace string) (map[string][]domain.Pod, error) {
	var (
		podLists = make(map[string][]domain.Pod)
		mu       sync.Mutex
		wg       sync.WaitGroup
		firstErr error
	)

	for key, repo := range pi.PodRepo {
		wg.Go(func() {
			pods, err := repo.GetAll(ctx, namespace)

			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				if firstErr == nil {
					firstErr = fmt.Errorf("context %s: %w", key, err)
				}
				fmt.Printf("Error fetching pods from context %s: %v\n", key, err)
				return
			}
			podLists[key] = pods
		})
	}
	wg.Wait()
	return podLists, firstErr
}

func (pi *podInteractor) GetFromManyContext(ctx context.Context, namespaces []string, contexts []string) (map[string][]domain.Pod, error) {
	var (
		podLists = make(map[string][]domain.Pod)
		mu       sync.Mutex
		wg       sync.WaitGroup
		firstErr error
	)

	for _, clusterCtx := range contexts {
		repo, ok := pi.PodRepo[clusterCtx]
		if !ok {
			continue
		}
		for _, ns := range namespaces {
			wg.Go(func() {
				pods, err := repo.GetAll(ctx, ns)

				mu.Lock()
				defer mu.Unlock()
				if err != nil {
					if firstErr == nil {
						firstErr = fmt.Errorf("context %s, namespace %s: %w", clusterCtx, ns, err)
					}
					fmt.Printf("Error fetching pods from context %s, namespace %s: %v\n", clusterCtx, ns, err)
					return
				}
				podLists[clusterCtx] = append(podLists[clusterCtx], pods...)
			})
		}
	}
	wg.Wait()
	return podLists, firstErr
}

func (pi *podInteractor) GetAllOneContext(ctx context.Context, namespace string, context string) ([]domain.Pod, error) {
	repo, ok := pi.PodRepo[context]
	if !ok {
		return nil, fmt.Errorf("no gateway found for context: %s", context)
	}
	pods, err := repo.GetAll(ctx, namespace)
	if err != nil {
		fmt.Printf("Error fetching pods from context %s: %v\n", context, err)
		return nil, err
	}
	return pods, nil
}

func (pi *podInteractor) PortForward(ctx context.Context, podName, namespace, context string, ports []string, stopChan <-chan struct{}, readyChan chan struct{}) (*bytes.Buffer, *bytes.Buffer, error) {
	gateway := pi.PodRepo[context]
	if gateway == nil {
		return nil, nil, fmt.Errorf("no gateway found for context: %s", context)
	}
	return gateway.PortForward(namespace, podName, ports, stopChan, readyChan)
}

func (pi *podInteractor) GetYaml(ctx context.Context, namespace string, name string, context string) ([]byte, error) {
	return pi.PodRepo[context].GetYaml(ctx, namespace, name)
}
