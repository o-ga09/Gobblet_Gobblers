package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"main/config"

	"github.com/go-redis/redis/v8"
)

type KVS struct {
	Cli *redis.Client
}

func NewKVS(ctx context.Context, cfg *config.Config) (*KVS, error) {
	cli := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",cfg.RedisHost,cfg.RedisPort),
	})
	if err := cli.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return &KVS{Cli: cli}, nil
}

func (k *KVS) Save(ctx context.Context, key string, value int) error {
	id := int64(value)
	return k.Cli.Set(ctx,key,id,30*time.Minute).Err()
}

func (k *KVS) Load(ctx context.Context, key string) (int, error) {
	id, err := k.Cli.Get(ctx,key).Int64()
	if err != nil {
		return 0, fmt.Errorf("failed to get by %q: %w", key, errors.New("not found"))
	}
	return int(id), nil
}