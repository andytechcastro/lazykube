package controller

import (
	"bytes"
	"context"
	"errors"
	"io"
	"lazykube/internal/domain"
	"lazykube/internal/usecase"

	"k8s.io/client-go/tools/remotecommand"
)

type deploymentController struct {
	DeploymentInteractor usecase.DeploymentInteractor
	PodInteractor        usecase.PodInteractor
}

// NewDeploymentController return a controller
func NewDeploymentController(dinteractor usecase.DeploymentInteractor, pinteractor usecase.PodInteractor) ControllerResource {
	return &deploymentController{
		DeploymentInteractor: dinteractor,
		PodInteractor:        pinteractor,
	}
}

func (dC *deploymentController) PortForward(ctx context.Context, resourceName, namespace, context string, ports []string, stopChan <-chan struct{}, readyChan chan struct{}) (*bytes.Buffer, *bytes.Buffer, error) {
	return nil, nil, errors.New("port-forward not directly supported for deployments, use GetPods to select a pod first")
}

func (dC *deploymentController) GetPods(ctx context.Context, deploymentName, namespace, context string) ([]domain.Pod, error) {
	return dC.DeploymentInteractor.GetPods(ctx, deploymentName, namespace, context)
}

func (dC *deploymentController) Exec(ctx context.Context, deploymentName, namespace, context, command, containerName string, dryRun bool, options remotecommand.StreamOptions) error {
	return errors.New("exec not directly supported for deployments, use GetPods to select a pod first")
}

func (dC *deploymentController) GetAll(ctx context.Context, namespace string) (map[string][]map[string]string, error) {
	deploymentLists, err := dC.DeploymentInteractor.GetAll(ctx, namespace)
	if err != nil {
		return nil, err
	}

	return DeploymentListsToMaps(deploymentLists), nil
}

func (dC *deploymentController) GetFromManyContext(ctx context.Context, namespaces []string, contexts []string) (map[string][]map[string]string, error) {
	deployments, err := dC.DeploymentInteractor.GetFromManyContext(ctx, namespaces, contexts)
	if err != nil {
		return nil, err
	}
	return DeploymentListsToMaps(deployments), nil
}

func (dC *deploymentController) GetAllOneContext(ctx context.Context, namespace string, context string) ([]map[string]string, error) {
	deployments, err := dC.DeploymentInteractor.GetAllOneContext(ctx, namespace, context)
	if err != nil {
		return nil, err
	}
	return DeploymentsToMaps(deployments), nil
}

func (dC *deploymentController) GetYaml(ctx context.Context, namespace string, name string, context string) ([]byte, error) {
	return dC.DeploymentInteractor.GetYaml(ctx, namespace, name, context)
}

func (dC *deploymentController) GetLogs(ctx context.Context, resourceName, namespace, context, containerName string) (io.ReadCloser, error) {
	return nil, errors.New("logs not directly supported for deployments, use GetPods to select a pod first")
}
