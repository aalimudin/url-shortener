package main

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

const expiryTime = 10

type urlData struct {
	ShortURL  string `json:"shortURL"`
	ActualURL string `json:"actualURL"`
}

type redisRepositoryInterface interface {
	Set(ctx context.Context, key string, data interface{}) error
	Get(ctx context.Context, key string) (*urlData, error)
}

type redisRepository struct {
	client *redis.Client
}

func newRedisRepository(address, password string, db int) redisRepositoryInterface {
	cli := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})

	return redisRepository{
		client: cli,
	}
}

func (r redisRepository) Set(ctx context.Context, key string, data interface{}) error {
	res, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return err
	}

	if res == 1 {
		return errors.New("shortened url already exists")
	}
	return r.client.Set(ctx, key, data, expiryTime*time.Minute).Err()
}

func (r redisRepository) Get(ctx context.Context, key string) (*urlData, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	urlDataRes := urlData{}

	err = json.Unmarshal([]byte(val), &urlDataRes)
	if err != nil {
		return nil, err
	}

	return &urlDataRes, nil
}
