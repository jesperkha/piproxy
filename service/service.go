package service

import (
	"encoding/json"
	"io"
	"os"
)

type Service struct {
	Name     string `json:"name"`
	Endpoint string `json:"endpoint"`
	Url      string `json:"url"`
}

func Load(path string) ([]Service, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var services []Service
	err = json.Unmarshal(b, &services)
	return services, nil
}
