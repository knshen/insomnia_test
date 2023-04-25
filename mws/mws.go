package mws

import (
	"context"
	"net/http"
)

type RequestHandler func(ctx context.Context, req *http.Request, writer http.ResponseWriter) (interface{}, error)

type Middleware func(RequestHandler) RequestHandler

func Build(mws ...Middleware) Middleware {
	return func(next RequestHandler) RequestHandler {
		for i := len(mws) - 1; i >= 0; i-- {
			next = mws[i](next)
		}
		return next
	}
}
