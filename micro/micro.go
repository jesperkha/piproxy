package micro

import (
	"context"
	"log"
	"net/http"

	"github.com/jesperkha/notifier"
)

type MicroService struct {
	mux  *http.ServeMux
	port string
	name string
}

func New(name string, port string) *MicroService {
	return &MicroService{
		mux:  http.NewServeMux(),
		port: port,
		name: name,
	}
}

func (m *MicroService) Handle(path string, f http.HandlerFunc) {
	m.mux.HandleFunc(path, f)
}

func (m *MicroService) ListenAndServe(notif *notifier.Notifier) error {
	s := &http.Server{
		Addr:    m.port,
		Handler: m.mux,
	}

	done, finish := notif.Register()

	go func() {
		<-done
		s.Shutdown(context.Background())
		log.Printf("%s: shutdown complete", m.name)
		finish()
	}()

	log.Printf("%s: listening at port %s", m.name, m.port)
	return s.ListenAndServe()
}
