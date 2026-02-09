package middleware

import (
	"context"
	"net/http"
	"net/http/httputil"
)

func NewMiddleware(proxy *httputil.ReverseProxy) *httputil.ReverseProxy {
	mwProxy := ContentIDHeader(proxy, "contentID")
	return mwProxy
}

func ContentIDHeader(proxy *httputil.ReverseProxy, contentID string) *httputil.ReverseProxy {
	defaultDirector := proxy.Director
	newDirector := func(r *http.Request) {
		defaultDirector(r)
		head := r.Header.Get("X-Content-ID")
		if head == "" {
			head = contentID
			r.Header.Add("X-Content-ID", head)
		}
		ctx := context.WithValue(r.Context(), "contentID", head)
		*r = *r.WithContext(ctx)
	}
	proxy.Director = newDirector
	proxy.ModifyResponse = func(r *http.Response) error {
		head, _ := r.Request.Context().Value("contentID").(string)
		r.Header.Set("X-Content-ID", head)
		return nil
	}
	return proxy
}
