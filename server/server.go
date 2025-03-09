package server

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"sync"

	"maps"

	"github.com/jesperkha/piproxy/config"
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

func (s *Server) ListenAndServe(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)

	server := &http.Server{
		Addr:    s.config.Port,
		Handler: s.mux,
	}

	go func() {
		<-ctx.Done()
		if err := server.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}

		log.Println("server shutdown")
		wg.Done()
	}()

	s.mux.HandleFunc("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
	}))

	log.Printf("listening at localhost:%s", s.config.Port)
	server.ListenAndServe()
}

func (s *Server) register(path string, host string) {
	localUrl, err := url.Parse(host)
	if err != nil {
		log.Fatal(err)
	}

	s.mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		redirectTo(r.URL, localUrl)

		proxy := http.DefaultTransport
		res, err := proxy.RoundTrip(r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		defer res.Body.Close()

		maps.Copy(w.Header(), res.Header)
		w.WriteHeader(res.StatusCode)

		if err := res.Write(w); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	})
}
