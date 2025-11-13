package registries

import (
	"lazykube/controllers"
	"lazykube/gateways"
	interGate "lazykube/usecases/gateways"
	"lazykube/usecases/interactors"
)

func (r *registry) NewDeploymentController() controllers.ControllerResource {
	deployGates := map[string]interGate.DeploymentResourceGateway{}
	for key, client := range r.clients {
		config := r.configs[key]
		deployGate := gateways.NewDeploymentGateway(client, config, key)
		deployGates[key] = deployGate
	}
	dInteractor := interactors.NewDeploymentInteractor(deployGates)

	podGates := map[string]interGate.PodResourceGateway{}
	for key, client := range r.clients {
		config := r.configs[key]
		podGate := gateways.NewPodGateway(client, config, key)
		podGates[key] = podGate
	}
	pInteractor := interactors.NewPodInteractor(podGates)

	return controllers.NewDeploymentController(dInteractor, pInteractor)
}
