package test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/jesperkha/piproxy/config"
	"github.com/jesperkha/piproxy/server"
	"github.com/jesperkha/piproxy/service"
)

func createServices() (services []service.Service, servers []*httptest.Server) {
	names := []string{
		"one", "two", "three",
	}

	for _, name := range names {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(name))
		}))

		servers = append(servers, server)
		services = append(services, service.Service{
			Name:     name,
			Endpoint: "/" + name,
			Url:      server.URL,
		})
	}

	return services, servers
}

func TestProxy(t *testing.T) {
	config := config.Config{
		Port: ":8080",
	}

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	services, servers := createServices()

	server := server.New(config)
	if err := server.RegisterServices(services); err != nil {
		t.Fatal(err)
	}

	go server.ListenAndServe(ctx, &wg)

	for _, s := range services {
		res, err := http.Get(fmt.Sprintf("http://localhost%s%s", config.Port, s.Endpoint))
		if err != nil {
			t.Error(err)
		}

		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("non-ok response from %s", s.Name)
		}

		// body, err := io.ReadAll(res.Body)
		// if err != nil {
		// 	t.Error(err)
		// }

		// if string(body) != s.Name {
		// 	t.Errorf("expected response '%s', got '%s'", s.Name, string(body))
		// }
	}

	cancel()
	wg.Wait()

	for _, s := range servers {
		s.Close()
	}
}
