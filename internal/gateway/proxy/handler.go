package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type FinalURL struct {
	scheme string
	host   string
	path   string
}

func HandleProxy(proxy *httputil.ReverseProxy) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	}
}

func NewProxy(targetHost string) (*httputil.ReverseProxy, error) {
	if !strings.Contains(targetHost, "://") {
		targetHost = "http://" + targetHost
	}

	u, err := url.Parse(targetHost)
	if err != nil {
		return nil, fmt.Errorf("could not parse target host: %w", err)
	}

	prxURL := FinalURL{
		scheme: u.Scheme,
		host:   u.Host,
		path:   u.Path,
	}

	newProxy := httputil.NewSingleHostReverseProxy(prxURL.URL())
	defaultDirector := newProxy.Director
	newDirector := func(r *http.Request) {
		defaultDirector(r)
		if r.URL.Path == "/proxy" {
			r.URL.Path = "/"
			return
		}
		if strings.HasPrefix(r.URL.Path, "/proxy/") {
			r.URL.Path = strings.TrimPrefix(r.URL.Path, "/proxy")
		}
	}
	newProxy.Director = newDirector
	return newProxy, nil
}

func (u FinalURL) URL() *url.URL {
	return &url.URL{
		Scheme: u.scheme,
		Host:   u.host,
		Path:   u.path,
	}
}

func (u FinalURL) String() string {
	return u.URL().String()
}
