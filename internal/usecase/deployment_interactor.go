package usecase

import (
	"context"
	"fmt"
	"lazykube/internal/domain"
	"lazykube/internal/usecase/port"
	"sync"
)

type deploymentInteractor struct {
	DeploymentRepo map[string]port.DeploymentResourceGateway
}

// DeploymentInteractor is an interface for connect to deployment interactor
type DeploymentInteractor interface {
	GetAll(context.Context, string) (map[string][]domain.Deployment, error)
	GetAllOneContext(context.Context, string, string) ([]domain.Deployment, error)
	GetFromManyContext(context.Context, []string, []string) (map[string][]domain.Deployment, error)
	GetYaml(context.Context, string, string, string) ([]byte, error)
	GetPods(ctx context.Context, deploymentName, namespace, context string) ([]domain.Pod, error)
}

// NewDeploymentInteractor return a new struct with deploymentInteractor
func NewDeploymentInteractor(
	deploymentRepo map[string]port.DeploymentResourceGateway,
) DeploymentInteractor {
	return &deploymentInteractor{
		DeploymentRepo: deploymentRepo,
	}
}

func (di *deploymentInteractor) GetPods(ctx context.Context, deploymentName, namespace, context string) ([]domain.Pod, error) {
	gateway := di.DeploymentRepo[context]
	if gateway == nil {
		return nil, fmt.Errorf("no gateway found for context: %s", context)
	}
	return gateway.GetPods(ctx, deploymentName, namespace)
}

func (di *deploymentInteractor) GetYaml(ctx context.Context, namespace string, name string, context string) ([]byte, error) {
	gateway, ok := di.DeploymentRepo[context]
	if !ok {
		return nil, fmt.Errorf("no gateway found for context: %s", context)
	}
	return gateway.GetYaml(ctx, namespace, name)
}

func (di *deploymentInteractor) GetAll(ctx context.Context, namespace string) (map[string][]domain.Deployment, error) {
	var (
		deploymentLists = make(map[string][]domain.Deployment)
		mu              sync.Mutex
		wg              sync.WaitGroup
		firstErr        error
	)

	for key, repo := range di.DeploymentRepo {
		wg.Go(func() {
			deployments, err := repo.GetAll(ctx, namespace)

			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				if firstErr == nil {
					firstErr = fmt.Errorf("context %s: %w", key, err)
				}
				fmt.Printf("Error fetching deployments from context %s: %v\n", key, err)
				return
			}
			deploymentLists[key] = deployments
		})
	}
	wg.Wait()
	return deploymentLists, firstErr
}

func (di *deploymentInteractor) GetFromManyContext(ctx context.Context, namespaces []string, contexts []string) (map[string][]domain.Deployment, error) {
	var (
		deploymentLists = make(map[string][]domain.Deployment)
		mu              sync.Mutex
		wg              sync.WaitGroup
		firstErr        error
	)

	for _, clusterCtx := range contexts {
		repo, ok := di.DeploymentRepo[clusterCtx]
		if !ok {
			continue
		}
		for _, ns := range namespaces {
			wg.Go(func() {
				deployments, err := repo.GetAll(ctx, ns)

				mu.Lock()
				defer mu.Unlock()
				if err != nil {
					if firstErr == nil {
						firstErr = fmt.Errorf("context %s, namespace %s: %w", clusterCtx, ns, err)
					}
					fmt.Printf("Error fetching deployments from context %s, namespace %s: %v\n", clusterCtx, ns, err)
					return
				}
				deploymentLists[clusterCtx] = append(deploymentLists[clusterCtx], deployments...)
			})
		}
	}
	wg.Wait()
	return deploymentLists, firstErr
}

func (di *deploymentInteractor) GetAllOneContext(ctx context.Context, namespace string, context string) ([]domain.Deployment, error) {
	repo, ok := di.DeploymentRepo[context]
	if !ok {
		return nil, fmt.Errorf("no gateway found for context: %s", context)
	}
	deployments, err := repo.GetAll(ctx, namespace)
	if err != nil {
		fmt.Printf("Error fetching deployments from context %s: %v\n", context, err)
		return nil, err
	}
	return deployments, nil
}
