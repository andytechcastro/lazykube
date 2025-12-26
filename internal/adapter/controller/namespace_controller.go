package controller

import (
	"context"
	"lazykube/internal/usecase"
)

type namespaceController struct {
	NamespaceInteractor usecase.NamespaceInteractor
}

func NewNamespaceController(interactor usecase.NamespaceInteractor) NamespaceController {
	return &namespaceController{
		NamespaceInteractor: interactor,
	}
}

func (nC namespaceController) GetAll(ctx context.Context, clusterContext string) ([]string, error) {
	return nC.NamespaceInteractor.GetAll(ctx, clusterContext)
}
