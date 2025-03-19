package main

import (
	"fmt"
	"log"
	"os"
	"syscall"

	"github.com/jesperkha/notifier"
	"github.com/jesperkha/piproxy/config"
	"github.com/jesperkha/piproxy/micro/services"
	"github.com/jesperkha/piproxy/micro/sysinfo"
	"github.com/jesperkha/piproxy/server"
	"github.com/jesperkha/piproxy/service"
)

func main() {
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

func url(host, port string) string {
	return fmt.Sprintf("http://%s%s", host, port)
}

func runMicroServices(s *server.Server, notif *notifier.Notifier, config config.Config) {
	sysinfoPort := ":5500"
	s.RegisterService("sysinfo", url(config.Host, sysinfoPort), "/sysinfo", func() {
		go sysinfo.Run(sysinfoPort, notif)
	})

	servicesPort := ":5501"
	s.RegisterService("services", url(config.Host, servicesPort), "/services", func() {
		go services.Run(servicesPort, notif, s)
	})
}
