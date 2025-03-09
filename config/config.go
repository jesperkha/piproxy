package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	ServiceFile string
}

func ensure(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("missing '%s' environment variable", k)
	}
	return v
}

func Load() Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	port := ensure("PORT")
	if port[0] != ':' {
		port = ":" + port
	}

	return Config{
		Port:        port,
		ServiceFile: ensure("SERVICE_PATH"),
	}
}
