package interactors

import (
	"context"
	"fmt"
	"lazykube/entities"
	"lazykube/usecases/gateways"
	"sync"
)

type deploymentInteractor struct {
	DeploymentRepo map[string]gateways.DeploymentResourceGateway
}

// DeploymentInteractor is an interface for connect to deployment interactor
type DeploymentInteractor interface {
	GetAll(string) (map[string][]entities.Deployment, error)
	GetAllOneContext(string, string) ([]entities.Deployment, error)
	GetFromManyContext([]string, []string) (map[string][]entities.Deployment, error)
	GetYaml(string, string, string) ([]byte, error)
	GetPods(ctx context.Context, deploymentName, namespace, context string) ([]entities.Pod, error)
}

// NewDeploymentInteractor return a new struct with deploymentInteractor
func NewDeploymentInteractor(
	deploymentRepo map[string]gateways.DeploymentResourceGateway,
) DeploymentInteractor {
	return &deploymentInteractor{
		DeploymentRepo: deploymentRepo,
	}
}

func (dI *deploymentInteractor) GetPods(ctx context.Context, deploymentName, namespace, context string) ([]entities.Pod, error) {
	gateway := dI.DeploymentRepo[context]
	if gateway == nil {
		return nil, fmt.Errorf("no gateway found for context: %s", context)
	}
	return gateway.GetPods(ctx, deploymentName, namespace)
}

func (dI *deploymentInteractor) GetAll(namespace string) (map[string][]entities.Deployment, error) {
	deploymentLists := map[string][]entities.Deployment{}
	wg := sync.WaitGroup{}
	var err error
	for key, repo := range dI.DeploymentRepo {
		wg.Add(1)
		go func(namespace string, key string, repo gateways.DeploymentResourceGateway) {
			deploymentLists[key], err = repo.GetAll(namespace)
			if err != nil {
				fmt.Println(err)
			}
			defer wg.Done()
		}(namespace, key, repo)
	}
	wg.Wait()
	return deploymentLists, nil
}

func (dI *deploymentInteractor) GetFromManyContext(namespaces []string, contexts []string) (map[string][]entities.Deployment, error) {
	deploymentLists := map[string][]entities.Deployment{}
	wg := sync.WaitGroup{}
	for _, context := range contexts {
		for _, namespace := range namespaces {
			wg.Add(1)
			repo := dI.DeploymentRepo[context]
			go func(namespace string, context string, repo gateways.DeploymentResourceGateway) {
				deploymentNamespaceLists, err := repo.GetAll(namespace)
				if err != nil {
					fmt.Println(err)
				}
				deploymentLists[context] = append(deploymentLists[context], deploymentNamespaceLists...)
				defer wg.Done()
			}(namespace, context, repo)
		}
	}
	wg.Wait()
	return deploymentLists, nil
}

func (dI *deploymentInteractor) GetAllOneContext(namespace string, context string) ([]entities.Deployment, error) {
	deploymentLists := []entities.Deployment{}
	var err error
	deploymentLists, err = dI.DeploymentRepo[context].GetAll(namespace)
	if err != nil {
		fmt.Println(err)
	}
	return deploymentLists, nil
}

func (dI *deploymentInteractor) GetYaml(namespace string, name string, context string) ([]byte, error) {
	return dI.DeploymentRepo[context].GetYaml(namespace, name)
}
