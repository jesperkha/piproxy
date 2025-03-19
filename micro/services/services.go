package services

import (
	"encoding/json"
	"net/http"

	"github.com/jesperkha/notifier"
	"github.com/jesperkha/piproxy/micro"
	"github.com/jesperkha/piproxy/server"
)

func Run(port string, notif *notifier.Notifier, server *server.Server) {
	m := micro.New("services", port)

	m.Handle("/", func(w http.ResponseWriter, r *http.Request) {
		s := server.Services()
		b, err := json.MarshalIndent(s, "", "\t")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(b)
	})

	m.ListenAndServe(notif)
}
