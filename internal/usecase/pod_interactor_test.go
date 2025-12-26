package usecase_test

import (
	"bytes"
	"context"
	"io"
	"k8s.io/client-go/tools/remotecommand"
	"lazykube/internal/domain"
	"lazykube/internal/usecase"
	"lazykube/internal/usecase/port"
	"testing"
)

type mockPodGateway struct{}

func (m *mockPodGateway) GetAll(ctx context.Context, namespace string) ([]domain.Pod, error) {
	return []domain.Pod{{Name: "test-pod"}}, nil
}
func (m *mockPodGateway) GetByName(ctx context.Context, namespace string, name string) (*domain.Pod, error) {
	return &domain.Pod{Name: name}, nil
}
func (m *mockPodGateway) GetByLabels(ctx context.Context, namespace string, label map[string]string) ([]domain.Pod, error) {
	return []domain.Pod{}, nil
}
func (m *mockPodGateway) GetYaml(ctx context.Context, namespace string, name string) ([]byte, error) {
	return []byte("yaml"), nil
}
func (m *mockPodGateway) Exec(ctx context.Context, podName, namespace, command, containerName string, dryRun bool, options remotecommand.StreamOptions) error {
	return nil
}
func (m *mockPodGateway) GetLogs(ctx context.Context, podName, namespace, containerName string) (io.ReadCloser, error) {
	return nil, nil
}
func (m *mockPodGateway) PortForward(namespace, podName string, ports []string, stopChan <-chan struct{}, readyChan chan struct{}) (*bytes.Buffer, *bytes.Buffer, error) {
	return nil, nil, nil
}

func TestPodInteractor_GetAll_Race(t *testing.T) {
	repos := make(map[string]port.PodResourceGateway)
	for i := 0; i < 10; i++ {
		repos[string(rune('a'+i))] = &mockPodGateway{}
	}

	pi := usecase.NewPodInteractor(repos)
	_, err := pi.GetAll(context.Background(), "default")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}
