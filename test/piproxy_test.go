package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
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

	services, servers := createServices()

	s := server.New(config)
	if err := s.RegisterServices(services); err != nil {
		t.Fatal(err)
	}

	notif := server.NewNotifier()
	go s.ListenAndServe(notif)

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

	notif.NotifyAndWait()

	for _, s := range servers {
		s.Close()
	}
}
