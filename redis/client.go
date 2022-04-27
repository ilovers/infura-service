package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type Client struct {
	redis *redis.Client
}

func NewClient(url string, pass string, db int) *Client {
	redisCli := redis.NewClient(&redis.Options{
		Addr:     url,
		Password: pass, // no password set
		DB:       db,   // use select DB
	})
	return &Client{
		redis: redisCli,
	}
}

func (c *Client) Get(key string) (string, error) {
	value, err := c.redis.Get(context.Background(), key).Result()
	return value, err
}
