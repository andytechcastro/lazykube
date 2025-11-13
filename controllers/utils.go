package controllers

import (
	"encoding/json"
	"fmt"
)

type ControllerResponse interface {
	map[string][]map[string]string | []map[string]string
}

func ResourceToData[T ControllerResponse](obj any) (T, error) {
	jsonBytes, err := json.Marshal(obj)
	var result T
	if err != nil {
		return result, fmt.Errorf("error al serializar a JSON: %w", err)
	}

	err = json.Unmarshal(jsonBytes, &result)
	if err != nil {
		return result, fmt.Errorf("error al deserializar a map: %w", err)
	}
	return result, nil
}
