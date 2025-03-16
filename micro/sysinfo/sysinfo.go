package sysinfo

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/jesperkha/notifier"
	"github.com/shirou/gopsutil/cpu"
)

func cpuInfo() (data string, err error) {
	info, err := cpu.Info()
	if err != nil {
		return data, err
	}

	for _, c := range info {
		return c.String(), err
	}

	return data, fmt.Errorf("sysinfo: failed to get CPU info")
}

func Run(host string, port string, notif *notifier.Notifier) {
	done, finish := notif.Register()
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		info, err := cpuInfo()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write([]byte(info))
	})

	s := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	go func() {
		<-done
		s.Shutdown(context.Background())
		log.Println("sysinfo: shutdown complete")
		finish()
	}()

	log.Printf("sysinfo: listening to %s%s", host, port)
	s.ListenAndServe()
}
