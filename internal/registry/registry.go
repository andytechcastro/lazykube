package registry

import (
	"lazykube/internal/adapter/controller"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type registry struct {
	clients map[string]kubernetes.Interface
	configs map[string]*rest.Config
}

// Registry registry for all the layers
type Registry interface {
	NewAppController() controller.AppController
	GetKubeConnection(context string) (kubernetes.Interface, *rest.Config)
}

// NewRegistry return a new registry
func NewRegistry(clients map[string]kubernetes.Interface, configs map[string]*rest.Config) Registry {
	return &registry{clients, configs}
}

func (r *registry) NewAppController() controller.AppController {
	return controller.AppController{
		Deployment: r.NewDeploymentController(),
		Pod:        r.NewPodController(),
		Namespace:  r.NewNamespaceController(),
	}
}

func (r *registry) GetKubeConnection(context string) (kubernetes.Interface, *rest.Config) {
	client := r.clients[context]
	config := r.configs[context]
	return client, config
}
