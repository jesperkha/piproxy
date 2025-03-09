package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"maps"

	"github.com/jesperkha/piproxy/config"
	"github.com/jesperkha/piproxy/service"
)

type Server struct {
	mux    *http.ServeMux
	config config.Config
}

func New(config config.Config) *Server {
	return &Server{
		config: config,
		mux:    http.NewServeMux(),
	}
}

func (s *Server) RegisterServices(services []service.Service) error {
	for _, serv := range services {
		serviceUrl, err := url.Parse(serv.Url)
		if err != nil {
			return err
		}

		if serv.Endpoint[0] != '/' {
			return fmt.Errorf("endpoint must start with '/': %s", serv.Endpoint)
		}

		s.register(serv.Endpoint, serviceUrl)
		log.Printf("registered service: %s", serv.Name)
	}

	return nil
}

func (s *Server) ListenAndServe(notif *Notifier) {
	done, finish := notif.Register()

	server := &http.Server{
		Addr:    s.config.Port,
		Handler: s.mux,
	}

	go func() {
		<-done
		if err := server.Shutdown(context.Background()); err != nil {
			log.Fatal(err)
		}

		log.Println("server shutdown")
		finish()
	}()

	log.Printf("listening at localhost%s", s.config.Port)
	server.ListenAndServe()
}

func (s *Server) register(path string, serviceUrl *url.URL) {
	s.mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		from := r.URL.String()
		redirectTo(r.URL, serviceUrl)

		log.Printf("new request: %s -> %s", from, r.URL.String())

		proxy := http.DefaultTransport
		res, err := proxy.RoundTrip(r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		defer res.Body.Close()

		if err := res.Write(w); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		maps.Copy(w.Header(), res.Header)
	})
}
