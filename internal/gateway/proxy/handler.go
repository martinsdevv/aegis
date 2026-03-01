package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/martinsdevv/aegis/internal/gateway/middleware"
)

func NewDynamicProxy() *httputil.ReverseProxy {
	return &httputil.ReverseProxy{
		Director: func(r *http.Request) {

			apiKey, ok := middleware.APIKeyFromContext(r.Context())
			if !ok || apiKey.UpstreamHost == "" {
				return
			}

			target := apiKey.UpstreamHost

			if !strings.Contains(target, "://") {
				target = "http://" + target
			}

			u, err := url.Parse(target)
			if err != nil {
				return
			}

			r.URL.Scheme = u.Scheme
			r.URL.Host = u.Host
			r.Host = u.Host

			if r.URL.Path == "/proxy" {
				r.URL.Path = "/"
			} else if strings.HasPrefix(r.URL.Path, "/proxy/") {
				r.URL.Path = strings.TrimPrefix(r.URL.Path, "/proxy")
			}

			ctx := middleware.SetUpstreamHost(r.Context(), u.Host)
			*r = *r.WithContext(ctx)
		},
	}
}

func HandleProxy(proxy *httputil.ReverseProxy) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	}
}
