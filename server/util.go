package server

import (
	"net/url"
	"strings"
)

func redirectTo(req *url.URL, to *url.URL) {
	req.Host = to.Host
	req.Scheme = to.Scheme
	req.Path = removeFirstPathSegment(req.Path)
}

func removeFirstPathSegment(path string) string {
	trimmed := strings.Trim(path, "/")
	split := strings.Split(trimmed, "/")

	if len(split) > 1 {
		return "/" + strings.Join(split[1:], "/")
	}

	return "/"
}
