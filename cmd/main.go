package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/jesperkha/piproxy/config"
	"github.com/jesperkha/piproxy/server"
	"github.com/jesperkha/piproxy/service"
)

func main() {
	config := config.Load()

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	s := server.New(config)
	services, err := service.Load(config.ServiceFile)
	if err != nil {
		log.Fatal(err)
	}

	if err := s.RegisterServices(services); err != nil {
		log.Fatal(err)
	}

	go s.ListenAndServe(ctx, &wg)

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)

	<-sigchan
	cancel()

	log.Println("piproxy shutting down")
	wg.Wait()
}
