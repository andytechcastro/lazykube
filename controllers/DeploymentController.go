package controllers

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"lazykube/entities"
	"lazykube/usecases/interactors"
	"log"

	"k8s.io/client-go/tools/remotecommand"
)

type deploymentController struct {
	DeploymentInteractor interactors.DeploymentInteractor
	PodInteractor        interactors.PodInteractor
}

// NewDeploymentController return a controller
func NewDeploymentController(dinteractor interactors.DeploymentInteractor, pinteractor interactors.PodInteractor) ControllerResource {
	return &deploymentController{
		DeploymentInteractor: dinteractor,
		PodInteractor:        pinteractor,
	}
}

func (dC *deploymentController) PortForward(ctx context.Context, resourceName, namespace, context string, ports []string, stopChan <-chan struct{}, readyChan chan struct{}) (*bytes.Buffer, *bytes.Buffer, error) {
	return nil, nil, errors.New("port-forward not directly supported for deployments, use GetPods to select a pod first")
}

func (dC *deploymentController) GetPods(ctx context.Context, deploymentName, namespace, context string) ([]entities.Pod, error) {
	return dC.DeploymentInteractor.GetPods(ctx, deploymentName, namespace, context)
}

func (dC *deploymentController) Exec(ctx context.Context, deploymentName, namespace, context, command, containerName string, dryRun bool, options remotecommand.StreamOptions) error {
	// TODO: Re-implement this with pod selection
	return nil
	// return dC.DeploymentInteractor.Exec(ctx, deploymentName, namespace, context, command, containerName, dryRun, options)
}

func (dC *deploymentController) GetAll(namespace string) (map[string][]map[string]string, error) {
	deploymentLists, err := dC.DeploymentInteractor.GetAll(namespace)
	if err != nil {
		fmt.Println(err)
		log.Panic()
	}

	return ResourceToData[map[string][]map[string]string](deploymentLists)
}

func (dC *deploymentController) GetFromManyContext(namespaces []string, contexts []string) (map[string][]map[string]string, error) {
	deployments, err := dC.DeploymentInteractor.GetFromManyContext(namespaces, contexts)
	if err != nil {
		return nil, err
	}
	return ResourceToData[map[string][]map[string]string](deployments)
}

func (dC *deploymentController) GetAllOneContext(namespace string, context string) ([]map[string]string, error) {
	deployments, err := dC.DeploymentInteractor.GetAllOneContext(namespace, context)
	if err != nil {
		return nil, err
	}
	return ResourceToData[[]map[string]string](deployments)
}

func (dC *deploymentController) GetYaml(namespace string, name string, context string) ([]byte, error) {
	return dC.DeploymentInteractor.GetYaml(namespace, name, context)
}

func (dC *deploymentController) GetLogs(ctx context.Context, resourceName, namespace, context, containerName string) (io.ReadCloser, error) {
	// TODO: Re-implement this with pod selection
	return nil, nil
	// return dC.DeploymentInteractor.GetLogs(ctx, resourceName, namespace, context, containerName)
}
