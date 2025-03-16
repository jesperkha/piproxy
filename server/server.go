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

func (s *Server) RegisterService(name string, url string, endpoint string, run func()) error {
	err := s.RegisterServices([]service.Service{
		{
			Name:     name,
			Url:      url,
			Endpoint: endpoint,
		},
	})

	if err != nil {
		return err
	}

	run()
	return nil
}

func (s *Server) RegisterServices(services []service.Service) error {
	for _, serv := range services {
		if err := s.register(serv); err != nil {
			return err
		}

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
			log.Println(err)
		}

		log.Println("server: shutdown complete")
		finish()
	}()

	log.Printf("server: listening at %s%s", s.config.Host, s.config.Port)
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

func (s *Server) register(serv service.Service) error {
	serviceUrl, err := url.Parse(serv.Url)
	if err != nil {
		return err
	}

	if serv.Endpoint[0] != '/' {
		return fmt.Errorf("server: endpoint must start with '/': %s", serv.Endpoint)
	}

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

	return nil
}
