package k8s_test

import (
	"context"
	"fmt"
	"lazykube/internal/infrastructure/k8s"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestCreatePod(t *testing.T) {
	client := fake.NewSimpleClientset(&v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "influxdb-v2",
			Namespace:   "tools",
			Annotations: map[string]string{},
		},
		Status: v1.PodStatus{
			Phase: "Running",
			ContainerStatuses: []v1.ContainerStatus{
				{
					Name:         "istio",
					RestartCount: 0,
					Image:        "gcr.io/coverwallet-sre/istio:2",
				},
				{
					Name:         "operator",
					RestartCount: 2,
					Image:        "gcr.io/coverwallet-sre/operator:2",
				},
			},
		},
	}, &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "chronograf",
			Namespace:   "default",
			Annotations: map[string]string{},
		},
		Status: v1.PodStatus{
			Phase: "Running",
			ContainerStatuses: []v1.ContainerStatus{
				{
					Name:         "istio",
					RestartCount: 0,
					Image:        "gcr.io/coverwallet-sre/istio:2",
				},
				{
					Name:         "operator",
					RestartCount: 0,
					Image:        "gcr.io/coverwallet-sre/app:2",
				},
			},
		},
	})
	operator := k8s.NewPodGateway(client, nil, "test-cluster")
	list, err := operator.GetAll(context.Background(), "tools")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(list)
}
