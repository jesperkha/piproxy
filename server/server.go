package server

import (
	"context"
	"log"
	"net/http"
	"sync"

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
