package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"maps"

	"github.com/jesperkha/notifier"
	"github.com/jesperkha/piproxy/config"
	"github.com/jesperkha/piproxy/service"
)

type Server struct {
	mux     *http.ServeMux
	config  config.Config
	handler http.Handler
}

func New(config config.Config) *Server {
	mux := http.NewServeMux()
	return &Server{
		config:  config,
		mux:     mux,
		handler: mux,
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

		s.register(serv, serviceUrl)
		log.Printf("server: registered service: %s for %s", serv.Name, serv.Endpoint)
	}

	return nil
}

func (s *Server) Middleware(middleware ...Middleware) {
	for _, m := range middleware {
		s.handler = m(s.handler)
	}
}

func (s *Server) ListenAndServe(notif *notifier.Notifier) {
	done, finish := notif.Register()

	server := &http.Server{
		Addr:    s.config.Port,
		Handler: s.handler,
	}

	go func() {
		<-done
		if err := server.Shutdown(context.Background()); err != nil {
			log.Fatal(err)
		}

		log.Println("server: shutdown complete")
		finish()
	}()

	log.Printf("server: listening at localhost%s", s.config.Port)
	server.ListenAndServe()
}

func (s *Server) register(serv service.Service, serviceUrl *url.URL) {
	s.mux.HandleFunc(serv.Endpoint, func(w http.ResponseWriter, r *http.Request) {
		redirectTo(r.URL, serviceUrl)

		proxy := http.DefaultTransport
		res, err := proxy.RoundTrip(r)
		if err != nil {
			log.Printf("server: '%s' failed to respond", serv.Name)
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
