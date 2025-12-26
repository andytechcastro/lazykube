package registry

import (
	"lazykube/internal/adapter/controller"
	"lazykube/internal/infrastructure/k8s"
	interGate "lazykube/internal/usecase/port"
	"lazykube/internal/usecase"
)

func (r *registry) NewNamespaceController() controller.NamespaceController {
	namesGates := map[string]interGate.NamespaceGateway{}
	for key, client := range r.clients {
		namesGate := k8s.NewNamespaceGateway(client, key)
		namesGates[key] = namesGate
	}

	return controller.NewNamespaceController(
		usecase.NewNamespaceInteractor(namesGates),
	)
}
