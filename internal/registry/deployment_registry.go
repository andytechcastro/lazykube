package registry

import (
	"lazykube/internal/adapter/controller"
	"lazykube/internal/infrastructure/k8s"
	interGate "lazykube/internal/usecase/port"
	"lazykube/internal/usecase"
)

func (r *registry) NewDeploymentController() controller.ControllerResource {
	deployGates := map[string]interGate.DeploymentResourceGateway{}
	for key, client := range r.clients {
		config := r.configs[key]
		deployGate := k8s.NewDeploymentGateway(client, config, key)
		deployGates[key] = deployGate
	}
	dInteractor := usecase.NewDeploymentInteractor(deployGates)

	podGates := map[string]interGate.PodResourceGateway{}
	for key, client := range r.clients {
		config := r.configs[key]
		podGate := k8s.NewPodGateway(client, config, key)
		podGates[key] = podGate
	}
	pInteractor := usecase.NewPodInteractor(podGates)

	return controller.NewDeploymentController(dInteractor, pInteractor)
}
