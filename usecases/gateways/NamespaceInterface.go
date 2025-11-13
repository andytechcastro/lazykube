package gateways

type NamespaceGateway interface {
	GetAll() ([]string, error)
}
