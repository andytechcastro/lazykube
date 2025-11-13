package main

import (
	"lazykube/infrastructure/config"
	"lazykube/infrastructure/datastore"
	"lazykube/infrastructure/tui"
	"lazykube/registries"
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

	registry := registries.NewRegistry(clients, configs)
	controllers := registry.NewAppController()
	clusters := maps.Keys(clients)
	tui.NewApp(clusters, controllers, conf)
}
