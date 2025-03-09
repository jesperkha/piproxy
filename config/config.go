package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
}

func Load() Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("missing 'PORT' environment variable")
	}

	if port[0] != ':' {
		port = ":" + port
	}

	return Config{
		Port: port,
	}
}
