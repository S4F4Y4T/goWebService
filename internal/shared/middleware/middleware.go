package middleware

import (
	"net/http"
)

type Middleware func(http.Handler) http.Handler

func Apply(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}

func With(h http.HandlerFunc, middlewares ...Middleware) http.Handler {
	return Apply(middlewares...)(h)
}
