package controllers

import "lazykube/usecases/interactors"

type namespaceController struct {
	NamespaceInteractor interactors.NamespaceInteractor
}

type NamespaceController interface {
	GetAll(string) ([]string, error)
}

func NewNamespaceController(interactor interactors.NamespaceInteractor) NamespaceController {
	return &namespaceController{
		NamespaceInteractor: interactor,
	}
}

func (nC namespaceController) GetAll(context string) ([]string, error) {
	return nC.NamespaceInteractor.GetAll(context)
}
