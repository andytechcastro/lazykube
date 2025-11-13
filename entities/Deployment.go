package entities

// Deployment the struct for the deployment information
type Deployment struct {
	Name              string            `json:"name,omitempty"`
	Namespace         string            `json:"namespace,omitempty"`
	Context           string            `json:"context,omitempty"`
	Replicas          int32             `json:"replicas,omitempty"`
	MatchLabels       map[string]string `json:"match_labels,omitempty"`
	AvailableReplicas int32             `json:"available_replicas,omitempty"`
	ReadyReplicas     int32             `json:"ready_replicas,omitempty"`
	UpdatedReplicas   int32             `json:"updated_replicas,omitempty"`
	PodList           []Pod             `json:"pod_list,omitempty"`
}
