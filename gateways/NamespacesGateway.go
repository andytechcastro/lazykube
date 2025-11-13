package gateways

import (
	"context"
	"lazykube/usecases/gateways"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type namespaceGateway struct {
	client  kubernetes.Interface
	context string
}

func NewNamespaceGateway(client kubernetes.Interface, cluster string) gateways.NamespaceGateway {
	return &namespaceGateway{
		client:  client,
		context: cluster,
	}
}

func (nG namespaceGateway) GetAll() ([]string, error) {
	namespaceList, err := nG.client.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return []string{}, err
	}
	namespaces := []string{}
	for _, namespace := range namespaceList.Items {
		namespaces = append(namespaces, namespace.GetName())
	}
	return namespaces, nil
}
