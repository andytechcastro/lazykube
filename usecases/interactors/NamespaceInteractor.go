package interactors

import "lazykube/usecases/gateways"

type namespaceInteractor struct {
	NamespaceGate map[string]gateways.NamespaceGateway
}

type NamespaceInteractor interface {
	GetAll(string) ([]string, error)
}

func NewNamespaceInteractor(namespaceGate map[string]gateways.NamespaceGateway) NamespaceInteractor {
	return &namespaceInteractor{
		NamespaceGate: namespaceGate,
	}
}

func (nI namespaceInteractor) GetAll(context string) ([]string, error) {
	if context == "all" {
		return []string{
			"all",
			"default",
		}, nil
	}
	return nI.NamespaceGate[context].GetAll()
}
