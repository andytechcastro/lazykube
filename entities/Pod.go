package entities

// Pod is the struct for save the resource info
type Pod struct {
	Name              string              `json:"name,omitempty"`
	Namespace         string              `json:"namespace,omitempty"`
	Context           string              `json:"context,omitempty"`
	State             string              `json:"state,omitempty"`
	ContainerStatuses []ContainerStatuses `json:"container_statuses,omitempty"`
}

// ContainerStatuses the struct for save container status
type ContainerStatuses struct {
	Name         string `json:"name,omitempty"`
	State        string `json:"state,omitempty"`
	RestartCount int32  `json:"restart_count,omitempty"`
	Image        string `json:"image,omitempty"`
}
