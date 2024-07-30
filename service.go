package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type shortenService struct {
	redisRepository redisRepositoryInterface
}

type shortenServiceInterface interface {
	ShortenURL(ctx context.Context, url, customSlug string) (string, error)
	GetShortenURL(ctx context.Context, slug string) (string, error)
}

func newShortenService(redisRepo redisRepositoryInterface) shortenServiceInterface {
	return &shortenService{
		redisRepository: redisRepo,
	}
}

func (s *shortenService) ShortenURL(ctx context.Context, url, customSlug string) (string, error) {
	key := generateShortKey(6)
	if customSlug != "" {
		key = customSlug
	}

	data := urlData{
		ShortURL:  fmt.Sprintf("localhost:8080/%s", key),
		ActualURL: url,
	}

	dataJSON, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	err = s.redisRepository.Set(ctx, key, dataJSON)
	if err != nil {
		return "", err
	}

	return data.ShortURL, nil
}

func (s *shortenService) GetShortenURL(ctx context.Context, slug string) (string, error) {
	urlData, err := s.redisRepository.Get(ctx, slug)
	if err != nil {
		if err == redis.Nil {
			return "", errors.New("url not found")
		}
		return "", err
	}
	return urlData.ActualURL, nil
}
