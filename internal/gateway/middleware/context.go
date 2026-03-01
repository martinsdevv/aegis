package middleware

import "context"

type ctxKeyUpstreamHost struct{}

func SetUpstreamHost(ctx context.Context, host string) context.Context {
	return context.WithValue(ctx, ctxKeyUpstreamHost{}, host)
}

func UpstreamHostFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(ctxKeyUpstreamHost{}).(string)
	return v, ok
}
