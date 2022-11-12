package db

import (
	"context"
	"github.com/go-redis/redis/v9"
	"log"
	"os"
)

var rdb *redis.Client

func init() {
	options, err := redis.ParseURL(os.Getenv("REDIS_ADDR"))
	options.Password = os.Getenv("REDIS_PASS")
	if err != nil {
		log.Fatal(err.Error())
	}
	rdb = redis.NewClient(options)
}

func SaveUrl(shortenedUrl string, targetUrl string) error {
	ctx := context.Background()
	err := rdb.Set(ctx, shortenedUrl, targetUrl, 0).Err()
	return err
}

func DeleteUrl(url string) error {
	ctx := context.Background()
	err := rdb.Del(ctx, url).Err()
	return err
}

func GetUrl(shortenedUrl string) (string, error) {
	ctx := context.Background()
	val, err := rdb.Get(ctx, shortenedUrl).Result()
	return val, err
}
