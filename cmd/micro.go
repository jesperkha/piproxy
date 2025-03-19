package cmd

import (
	"fmt"

	"github.com/jesperkha/notifier"
	"github.com/jesperkha/piproxy/config"
	"github.com/jesperkha/piproxy/micro/services"
	"github.com/jesperkha/piproxy/micro/sysinfo"
	"github.com/jesperkha/piproxy/server"
)

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
