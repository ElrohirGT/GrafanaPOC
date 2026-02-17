package main

import (
	"context"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/redis/go-redis/extra/redisotel/v9"
)

var rdb *redis.Client

func initRedis(ctx context.Context) error {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}

	rdb = redis.NewClient(&redis.Options{
		Addr:         addr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	})

	if err := redisotel.InstrumentTracing(rdb); err != nil {
		return err
	}
	if err := redisotel.InstrumentMetrics(rdb); err != nil {
		return err
	}

	// Ping to verify connection (non-fatal if Redis is optional for local dev)
	_ = rdb.Ping(ctx).Err()
	return nil
}

func closeRedis(ctx context.Context) error {
	if rdb != nil {
		return rdb.Close()
	}
	return nil
}

func getRedis() *redis.Client {
	return rdb
}
