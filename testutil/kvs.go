package testutil

import (
	"crypto/tls"
	"fmt"
	"os"
	"testing"

	"github.com/go-redis/redis/v8"
)

func OpenRedisForTest(t *testing.T) *redis.Client {
	t.Helper()

	env := os.Getenv("GAME_ENV")
	redis_host := os.Getenv("GAME_REDIS_HOST")
	redis_port := os.Getenv("GAME_REDIS_PORT")
	redis_password := os.Getenv("GAME_REDIS_PASSWORD")
	redis_TLS := os.Getenv("GAME_REDIS_TLS_SERVER_NAME")

	if _, defined := os.LookupEnv("CI"); defined {
		redis_port = "6379"
	}
	
	var opt *redis.Options
	switch (env) {
	case "dev":
		opt = &redis.Options{
			Addr: fmt.Sprintf("%s:%s",redis_host,redis_port),
			Password: redis_password,
			DB: 0,
		}
	case "prod":
		opt = &redis.Options{
			Addr: fmt.Sprintf("%s:%s",redis_host,redis_port),
			Password: redis_password,
			DB: 0,
			TLSConfig: &tls.Config{ServerName: redis_TLS},
		}

	}
	client := redis.NewClient(opt)
	return client
}