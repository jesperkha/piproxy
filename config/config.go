package config

import (
	"log"
	"os"

	"github.com/echo-webkom/cenv"
)

type Config struct {
	Port        string
	Host        string
	ServiceFile string
}

func Load() Config {
	if err := cenv.Load(); err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("PORT")
	if port[0] != ':' {
		port = ":" + port
	}

	return Config{
		Port:        port,
		Host:        os.Getenv("HOST"),
		ServiceFile: os.Getenv("SERVICE_PATH"),
	}
}
