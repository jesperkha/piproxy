package cmd

import (
	"log"
	"os"
	"syscall"

	"github.com/jesperkha/notifier"
	"github.com/jesperkha/piproxy/config"
	"github.com/jesperkha/piproxy/server"
	"github.com/jesperkha/piproxy/service"
)

func Run() {
	parseFlags()

	config := config.Load()
	notif := notifier.New()

	s := server.New(config)
	s.Middleware(server.Logger)

	servs, err := service.Load(config.ServiceFile)
	if err != nil {
		log.Fatal(err)
	}

	if err := s.RegisterServices(servs); err != nil {
		log.Fatal(err)
	}

	runMicroServices(s, notif, config)
	go s.ListenAndServe(notif)

	notif.NotifyOnSignal(os.Interrupt, syscall.SIGTERM)
	log.Println("piproxy: shutdown complete")
}
