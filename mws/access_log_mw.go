package mws

import (
	"context"
	"log"
	"net/http"
	"time"
)

func AccessLogMiddleware() Middleware {
	return func(next RequestHandler) RequestHandler {
		return func(ctx context.Context, req *http.Request, writer http.ResponseWriter) (interface{}, error) {
			log.Printf("start to handle req: %v", req.RequestURI)
			t0 := time.Now()
			res, err := next(ctx, req, writer)
			cost := time.Since(t0)
			log.Printf("finish handle req: %v, time cost: %v, resp header: %v", req.RequestURI, cost, writer.Header())
			return res, err
		}
	}
}
