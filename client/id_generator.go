package client

import (
	"context"
	"time"
)

type IIDGeneratClient interface {
	GenInt64ID(ctx context.Context, count int) ([]int64, error)
}

// suppose we have an id generate service that can safely gen unique int64 ids
type DefaultIDGenerateClient struct {
}

func (c *DefaultIDGenerateClient) GenInt64ID(ctx context.Context, count int) ([]int64, error) {
	// TODO call rpc of id generate service
	v := time.Now().Unix()
	res := make([]int64, 0, count)

	for i := 0; i < count; i++ {
		res = append(res, v+int64(i))
	}
	return res, nil
}
