// Storage
package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/dawidhermann/shortener-url/internal/config"
	"github.com/go-redis/redis/v9"
)

type KVStore struct {
	Rdb *redis.Client
}

func New(cfg config.StoreConfig) *KVStore {
	options, err := redis.ParseURL(cfg.Address)
	if err != nil {
		fmt.Println(os.Getenv("REDIS_ADDR"))
		log.Fatal(err.Error())
	}
	options.Password = cfg.Password
	rdb := redis.NewClient(options)
	return &KVStore{Rdb: rdb}
}

// Save url in key-value storage
func (store KVStore) SaveUrl(shortenedUrl string, targetUrl string) error {
	ctx := context.Background()
	err := store.Rdb.Set(ctx, shortenedUrl, targetUrl, 0).Err()
	return err
}

// Delete url from key-value storage
func (store KVStore) DeleteUrl(url string) error {
	ctx := context.Background()
	err := store.Rdb.Del(ctx, url).Err()
	return err
}
