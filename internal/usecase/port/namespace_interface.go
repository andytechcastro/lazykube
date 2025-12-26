package port

import "context"

type NamespaceGateway interface {
	GetAll(ctx context.Context) ([]string, error)
}
