package registries

import (
	"lazykube/controllers"
	"lazykube/gateways"
	interGate "lazykube/usecases/gateways"
	"lazykube/usecases/interactors"
)

func (r *registry) NewPodController() controllers.ControllerResource {
	podGates := map[string]interGate.PodResourceGateway{}
	for key, client := range r.clients {
		config := r.configs[key]
		podGate := gateways.NewPodGateway(client, config, key)
		podGates[key] = podGate
	}

	return controllers.NewPodController(
		interactors.NewPodInteractor(podGates),
	)
}
