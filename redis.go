package main

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var (
	rdb  *redis.Client
	rctx = context.Background()
)

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}
