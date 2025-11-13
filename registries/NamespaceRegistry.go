package registries

import (
	"lazykube/controllers"
	"lazykube/gateways"
	interGate "lazykube/usecases/gateways"
	"lazykube/usecases/interactors"
)

func (r *registry) NewNamespaceController() controllers.NamespaceController {
	namesGates := map[string]interGate.NamespaceGateway{}
	for key, client := range r.clients {
		namesGate := gateways.NewNamespaceGateway(client, key)
		namesGates[key] = namesGate
	}

	return controllers.NewNamespaceController(
		interactors.NewNamespaceInteractor(namesGates),
	)
}
