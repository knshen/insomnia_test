package mws

import (
	"context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	flag = []string{}
)

func TestMiddleware(t *testing.T) {
	chain := Build(Middleware1(), Middleware2())(Biz)
	ctx := context.Background()

	req := httptest.NewRequest("GET", "https://mock/mock", nil)
	w := httptest.NewRecorder()
	chain(ctx, req, w)

	a := assert.New(t)
	a.Equal(5, len(flag))
	a.Equal("start1", flag[0])
	a.Equal("start2", flag[1])
	a.Equal("biz", flag[2])
	a.Equal("end2", flag[3])
	a.Equal("end1", flag[4])
}

func Middleware1() Middleware {
	return func(next RequestHandler) RequestHandler {
		return func(ctx context.Context, req *http.Request, writer http.ResponseWriter) (interface{}, error) {
			flag = append(flag, "start1")
			res, err := next(ctx, req, writer)
			flag = append(flag, "end1")
			return res, err
		}
	}
}

func Middleware2() Middleware {
	return func(next RequestHandler) RequestHandler {
		return func(ctx context.Context, req *http.Request, writer http.ResponseWriter) (interface{}, error) {
			flag = append(flag, "start2")
			res, err := next(ctx, req, writer)
			flag = append(flag, "end2")
			return res, err
		}
	}
}

func Biz(ctx context.Context, req *http.Request, writer http.ResponseWriter) (interface{}, error) {
	flag = append(flag, "biz")
	return "ok", nil
}
