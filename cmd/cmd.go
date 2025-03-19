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
	flags := parseFlags()
	config := config.Load()

	if flags.LogFile {
		logFile, err := os.OpenFile(config.LogFile, os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal(err)
		}

		log.SetOutput(logFile)
		defer logFile.Close()
	}

	s := server.New(config)
	s.Middleware(server.Logger)

	servs, err := service.Load(config.ServiceFile)
	if err != nil {
		log.Fatal(err)
	}

	if err := s.RegisterServices(servs); err != nil {
		log.Fatal(err)
	}

	notif := notifier.New()

	runMicroServices(s, notif, config)
	go s.ListenAndServe(notif)

	notif.NotifyOnSignal(os.Interrupt, syscall.SIGTERM)
	log.Println("piproxy: shutdown complete")
}
