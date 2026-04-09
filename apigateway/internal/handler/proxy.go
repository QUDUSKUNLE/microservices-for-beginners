package handler

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func NewProxy(target string, prefix string) http.Handler {
	u, err := url.Parse(target)
	if err != nil {
		log.Fatalf("invalid proxy target URL %q: %v", target, err)
	}
	proxy := httputil.NewSingleHostReverseProxy(u)

	originalDirector := proxy.Director

	proxy.Director = func(req *http.Request) {
		originalDirector(req)

		path := strings.TrimPrefix(req.URL.Path, prefix)

		// ensure leading slash
		if path == "" {
			path = "/"
		} else if !strings.HasPrefix(path, "/") {
			path = "/" + path
		}

		req.URL.Path = path
		req.Host = u.Host
	}

	return proxy
}
