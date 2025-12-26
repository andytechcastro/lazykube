package usecase

import (
	"context"
	"lazykube/internal/usecase/port"
)

type namespaceInteractor struct {
	NamespaceGate map[string]port.NamespaceGateway
}

type NamespaceInteractor interface {
	GetAll(ctx context.Context, clusterContext string) ([]string, error)
}

func NewNamespaceInteractor(namespaceGate map[string]port.NamespaceGateway) NamespaceInteractor {
	return &namespaceInteractor{
		NamespaceGate: namespaceGate,
	}
}

func (nI namespaceInteractor) GetAll(ctx context.Context, clusterContext string) ([]string, error) {
	if clusterContext == "all" {
		return []string{
			"all",
			"default",
		}, nil
	}
	return nI.NamespaceGate[clusterContext].GetAll(ctx)
}
