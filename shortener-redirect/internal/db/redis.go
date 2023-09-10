package db

import (
	"context"
	"fmt"
	"log"

	"github.com/dawidhermann/shortener-redirect/internal/config"
	"github.com/go-redis/redis/v9"
)

type KVStore struct {
	Rdb *redis.Client
}

func New(cfg config.StoreConfig) *KVStore {
	options, err := redis.ParseURL(cfg.Address)
	if err != nil {
		fmt.Println(cfg.Address)
		log.Fatal(err.Error())
	}
	options.Password = cfg.Password
	return &KVStore{Rdb: redis.NewClient(options)}
}

func (store KVStore) GetUrl(shortenedUrl string) (string, error) {
	ctx := context.Background()
	val, err := store.Rdb.Get(ctx, shortenedUrl).Result()
	return val, err
}
