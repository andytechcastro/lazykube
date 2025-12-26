package main

import (
	"lazykube/internal/infrastructure/config"
	"lazykube/internal/infrastructure/datastore"
	"lazykube/internal/infrastructure/tui"
	"lazykube/internal/registry"
	"maps"
)

func main() {
	conf, err := config.ReadConfig()
	if err != nil {
		panic(1)
	}
	clients, configs, err := datastore.NewKubeConnections()
	if err != nil {
		panic(1)
	}

	registry := registry.NewRegistry(clients, configs)
	controllers := registry.NewAppController()
	clusters := maps.Keys(clients)
	tui.NewApp(clusters, controllers, conf)
}
