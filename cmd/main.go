package main

import (
	"log"
	"os"
	"syscall"

	"github.com/jesperkha/notifier"
	"github.com/jesperkha/piproxy/config"
	"github.com/jesperkha/piproxy/server"
	"github.com/jesperkha/piproxy/service"
)

func main() {
	config := config.Load()

	services, err := service.Load(config.ServiceFile)
	if err != nil {
		log.Fatal(err)
	}

	s := server.New(config)
	if err := s.RegisterServices(services); err != nil {
		log.Fatal(err)
	}

	notif := notifier.New()
	go s.ListenAndServe(notif)

	notif.NotifyOnSignal(os.Interrupt, syscall.SIGTERM)
	log.Println("piproxy shutting down")
}
