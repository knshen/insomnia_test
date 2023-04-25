package mws

import (
	"context"
	"net/http"
)

func UserTokenCheckerMiddleware() Middleware {
	return func(next RequestHandler) RequestHandler {
		return func(ctx context.Context, req *http.Request, writer http.ResponseWriter) (interface{}, error) {
			// check if session token is valid
			if !isValid(req.Header.Get("user_token")) {
				writer.WriteHeader(401)
				return nil, nil
			}
			return next(ctx, req, writer)
		}
	}
}

func isValid(token string) bool {
	// TODO
	return token != ""
}
