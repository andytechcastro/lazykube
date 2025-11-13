package datastore

import (
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

// NewKubeConnections return kubernetes clients and configs
func NewKubeConnections() (map[string]kubernetes.Interface, map[string]*rest.Config, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, nil, err
	}
	kubeConfigPath := filepath.Join(userHomeDir, ".kube", "config")

	name := getAllContext(kubeConfigPath)

	clients := map[string]kubernetes.Interface{}
	configs := map[string]*rest.Config{}

	for key := range name {
		kubeConfig, _ := buildConfigFromFlags(key, kubeConfigPath)

		client, err := kubernetes.NewForConfig(kubeConfig)
		if err != nil {
			return nil, nil, err
		}
		clients[key] = client
		configs[key] = kubeConfig
	}
	return clients, configs, nil
}

func buildConfigFromFlags(context, kubeconfigPath string) (*rest.Config, error) {
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath},
		&clientcmd.ConfigOverrides{
			CurrentContext: context,
		}).ClientConfig()
}

func getAllContext(pathToKubeConfig string) map[string]*api.Context {
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: pathToKubeConfig},
		&clientcmd.ConfigOverrides{
			CurrentContext: "",
		}).RawConfig()

	if err != nil {
		fmt.Println(err)
	}

	return config.Contexts
}
