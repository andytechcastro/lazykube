package gateways

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
)

// PortForwardGateway handles port forwarding to a pod.
type PortForwardGateway struct {
	config *rest.Config
}

// NewPortForwardGateway creates a new PortForwardGateway.
func NewPortForwardGateway(config *rest.Config) *PortForwardGateway {
	return &PortForwardGateway{
		config: config,
	}
}

// Forward establishes a port forwarding connection to a pod.
// It runs until the stopChan is closed.
func (g *PortForwardGateway) Forward(namespace, podName string, ports []string, stopChan <-chan struct{}, readyChan chan struct{}) (*bytes.Buffer, *bytes.Buffer, error) {
	// Build the URL for the port forward request
	path := fmt.Sprintf("/api/v1/namespaces/%s/pods/%s/portforward", namespace, podName)
	hostIP := strings.TrimLeft(g.config.Host, "htps:/")
	serverURL, err := url.Parse(fmt.Sprintf("https://%s%s", hostIP, path))
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing URL: %w", err)
	}

	// Create a SPDY transport
	transport, upgrader, err := spdy.RoundTripperFor(g.config)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating round tripper: %w", err)
	}

	// Create the dialer
	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, http.MethodPost, serverURL)

	// Create streams for output
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	// Create the port forwarder
	forwarder, err := portforward.New(dialer, ports, stopChan, readyChan, stdout, stderr)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating port forwarder: %w", err)
	}

	// Run the port forwarder in a goroutine
	go func() {
		// This will block until the stopChan is closed or an error occurs.
		if err = forwarder.ForwardPorts(); err != nil {
			// It's common for this to error out when the connection is closed,
			// so we might not want to log this as a fatal error.
			// For now, we can just print it to stderr buffer.
			fmt.Fprintf(stderr, "Error forwarding ports: %v\n", err)
		}
	}()

	return stdout, stderr, nil
}
