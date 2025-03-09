package server

import (
	"net/http"
	"testing"
)

func TestRemoveFirstPathSegment(t *testing.T) {
	cases := map[string]string{
		"":              "/",
		"/":             "/",
		"/users/john":   "/john",
		"/foo/bar/faz/": "/bar/faz",
		"bap/bow":       "/bow",
		"foo/bar/":      "/bar",
	}

	for in, expect := range cases {
		got := removeFirstPathSegment(in)
		if got != expect {
			t.Errorf("expected %s, got %s", expect, got)
		}
	}
}

func requestDeepEqual(r1, r2 *http.Request) bool {
	for k := range r1.Header {
		if r2.Header.Get(k) != r1.Header.Get(k) {
			return false
		}
	}

	url := r1.Method == r2.Method && r1.URL.String() == r2.URL.String()
	headerSize := len(r1.Header) == len(r2.Header)
	return url && headerSize
}

func TestRedirectTo(t *testing.T) {
	expect, err := http.NewRequest("GET", "http://localhost:1234/endpoint", nil)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", "http://host.com/service/endpoint", nil)
	if err != nil {
		t.Fatal(err)
	}

	redirectTo(req.URL, expect.URL)

	if !requestDeepEqual(req, expect) {
		t.Error("expected equal reqests")
	}
}
