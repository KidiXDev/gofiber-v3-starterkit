package redis

import (
	"context"
	"os"

	"gofiber-starterkit/pkg/utils"

	"github.com/gofiber/storage/redis/v3"
	redigo "github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type RedisClient struct {
	Client  *redigo.Client
	Storage *redis.Storage
}

func New() *RedisClient {
	rdb := redigo.NewClient(&redigo.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Username: os.Getenv("REDIS_USERNAME"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       utils.ParseIntEnv("REDIS_DB", 0),
		PoolSize: 20,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to Redis")
	}

	storage := redisStorageClient(rdb)

	log.Debug().Msg("Connected to Redis successfully")

	return &RedisClient{
		Client:  rdb,
		Storage: storage,
	}
}

func redisStorageClient(cli *redigo.Client) *redis.Storage {
	return redis.NewFromConnection(cli)
}
