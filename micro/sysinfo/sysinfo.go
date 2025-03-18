package sysinfo

import (
	"fmt"
	"net/http"

	"github.com/jesperkha/notifier"
	"github.com/jesperkha/piproxy/micro"
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

func Run(port string, notif *notifier.Notifier) {
	m := micro.New("sysinfo", port)

	m.Handle("/", func(w http.ResponseWriter, r *http.Request) {
		info, err := cpuInfo()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write([]byte(info))
	})

	m.ListenAndServe(notif)
}
