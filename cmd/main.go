package main

import (
	"fmt"
	"log"
	"os"
	"syscall"

	"github.com/jesperkha/notifier"
	"github.com/jesperkha/piproxy/config"
	"github.com/jesperkha/piproxy/micro/sysinfo"
	"github.com/jesperkha/piproxy/server"
	"github.com/jesperkha/piproxy/service"
)

func main() {
	config := config.Load()
	notif := notifier.New()

	services, err := service.Load(config.ServiceFile)
	if err != nil {
		log.Fatal(err)
	}

	s := server.New(config)
	s.Middleware(server.Logger)

	if err := s.RegisterServices(services); err != nil {
		log.Fatal(err)
	}

	s.RegisterService("sysinfo", fmt.Sprintf("http://%s:5500", config.Host), "/sysinfo", func() {
		go sysinfo.Run(":5500", notif)
	})

	go s.ListenAndServe(notif)

	notif.NotifyOnSignal(os.Interrupt, syscall.SIGTERM)
	log.Println("piproxy: shutdown complete")
}
