package entities

import "fmt"

type ContainerSelectionError struct {
	Containers []string
}

func (e *ContainerSelectionError) Error() string {
	return fmt.Sprintf("multiple containers found, please select one: %v", e.Containers)
}
