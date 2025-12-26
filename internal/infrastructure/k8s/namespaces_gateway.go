package k8s

import (
	"context"
	"fmt"
	"lazykube/internal/usecase/port"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type namespaceGateway struct {
	client  kubernetes.Interface
	context string
}

func NewNamespaceGateway(client kubernetes.Interface, cluster string) port.NamespaceGateway {
	return &namespaceGateway{
		client:  client,
		context: cluster,
	}
}

func (nG namespaceGateway) GetAll(ctx context.Context) ([]string, error) {
	namespaceList, err := nG.client.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return []string{}, fmt.Errorf("failed to list namespaces: %w", err)
	}
	namespaces := []string{}
	for _, namespace := range namespaceList.Items {
		namespaces = append(namespaces, namespace.GetName())
	}
	return namespaces, nil
}
