package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jesperkha/piproxy/config"
	"github.com/jesperkha/piproxy/server"
	"github.com/jesperkha/piproxy/service"
)

func main() {
	config := config.Load()

	s := server.New(config)
	services, err := service.Load(config.ServiceFile)
	if err != nil {
		log.Fatal(err)
	}

	if err := s.RegisterServices(services); err != nil {
		log.Fatal(err)
	}

	notif := server.NewNotifier()
	go s.ListenAndServe(notif)

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)

	<-sigchan
	notif.NofifyAndWait()

	log.Println("piproxy shutting down")
}
