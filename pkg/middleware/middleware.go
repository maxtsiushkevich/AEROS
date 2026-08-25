package middleware

import (
	"net/http"
)

type (
	middleware      func(http.HandlerFunc) http.HandlerFunc
	MiddlewareGroup []middleware
)

func (mg MiddlewareGroup) Apply(h http.HandlerFunc) http.HandlerFunc {
	for i := len(mg) - 1; i >= 0; i-- {
		h = mg[i](h)
	}

	return h
}
