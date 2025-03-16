package server

import (
	"context"
	"fmt"
	"io"
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

func (s *Server) handle(path string, f func(w http.ResponseWriter, r *http.Request) (int, error)) {
	s.mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		status, err := f(w, r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(status)
		}
	})
}

func (s *Server) register(serv service.Service, serviceUrl *url.URL) {
	s.handle(serv.Endpoint, func(w http.ResponseWriter, r *http.Request) (int, error) {
		redirectTo(r.URL, serviceUrl)
		serverError := http.StatusInternalServerError

		res, err := http.DefaultTransport.RoundTrip(r)
		if err != nil {
			return serverError, fmt.Errorf("server: '%s' failed to respond", serv.Name)
		}

		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return serverError, fmt.Errorf("server: failed to read request body")
		}

		if _, err := w.Write(body); err != nil {
			return serverError, fmt.Errorf("server: failed to write request body")
		}

		maps.Copy(w.Header(), res.Header)
		return http.StatusOK, nil
	})
}
