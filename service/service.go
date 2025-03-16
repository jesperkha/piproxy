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

func Load(paths ...string) (services []Service, err error) {
	for _, path := range paths {
		file, err := os.Open(path)
		defer file.Close()

		if err != nil {
			return nil, err
		}

		b, err := io.ReadAll(file)
		if err != nil {
			return nil, err
		}

		var s []Service
		err = json.Unmarshal(b, &s)
		services = append(services, s...)
	}

	return services, nil
}
