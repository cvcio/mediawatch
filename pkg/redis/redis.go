package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v9"
)

// RedisClient
type RedisClient struct {
	Client *redis.Client
	ctx    context.Context
}

// RedisClient
func NewRedisClient(ctx context.Context, url string, pass string) (*RedisClient, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:            url,
		Password:        pass,
		MaxRetries:      3,
		ConnMaxIdleTime: time.Minute * 5,
	})

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("Redis Ping/Pong failed with error: %s", err.Error())
	}

	return &RedisClient{
		Client: rdb,
		ctx:    ctx,
	}, nil
}

// Close closes the redis connection
func (rdb *RedisClient) Close() error {
	return rdb.Client.Close()
}

// Publish
func (rdb *RedisClient) Publish(key string, value string) error {
	return rdb.Client.Publish(rdb.ctx, key, value).Err()
}

// Subscribe to redis channel and return a channel of messages
func (rdb *RedisClient) Subscribe(ctx context.Context, key string) (*redis.PubSub, chan []byte, error) {
	pubsub := rdb.Client.PSubscribe(ctx, key)
	msg := make(chan []byte)
	go func(channel <-chan *redis.Message) {
		for m := range channel {
			msg <- []byte(m.Payload)
		}
	}(pubsub.Channel())
	return pubsub, msg, nil
}

// Get
func (rdb *RedisClient) Get(key string) (string, error) {
	val, err := rdb.Client.Get(rdb.ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key %s does not exists", key)
	} else if err != nil {
		return "", err
	}
	return val, nil
}

// Set
func (rdb *RedisClient) Set(key string, value string, ttl time.Duration) error {
	return rdb.Client.Set(rdb.ctx, key, value, ttl).Err()
}
