package mws

import (
	"code.sk.org/insomnia_test/consts"
	"context"
	"encoding/json"
	"net/http"
)

func RespHandlerMiddleware() Middleware {
	return func(next RequestHandler) RequestHandler {
		return func(ctx context.Context, req *http.Request, writer http.ResponseWriter) (interface{}, error) {
			res, err := next(ctx, req, writer)
			if err != nil {
				if code, ok := consts.Error2Code[err]; ok {
					writer.WriteHeader(int(code))
				} else {
					writer.WriteHeader(500)
				}
				writer.Write([]byte(err.Error()))
			} else {
				writer.WriteHeader(200)
				if str, ok := res.(string); ok {
					writer.Write([]byte(str))
				} else {
					body, _ := json.Marshal(res)
					writer.Write(body)
				}
			}

			return res, err
		}
	}
}
