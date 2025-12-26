package registry

import (
	"lazykube/internal/adapter/controller"
	"lazykube/internal/infrastructure/k8s"
	interGate "lazykube/internal/usecase/port"
	"lazykube/internal/usecase"
)

func (r *registry) NewPodController() controller.ControllerResource {
	podGates := map[string]interGate.PodResourceGateway{}
	for key, client := range r.clients {
		config := r.configs[key]
		podGate := k8s.NewPodGateway(client, config, key)
		podGates[key] = podGate
	}

	return controller.NewPodController(
		usecase.NewPodInteractor(podGates),
	)
}
