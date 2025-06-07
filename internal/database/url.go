package database

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func (s *service) StoreURL(ctx context.Context, shortCode, originalURL string) error {
	const oneWeek = 7 * 24 * time.Hour
	pipe := s.db.Pipeline()
	pipe.Set(ctx, "code:"+shortCode, originalURL, oneWeek)
	pipe.Set(ctx, "url:"+originalURL, shortCode, oneWeek)
	_, err := pipe.Exec(ctx)
	return err
}

func (s *service) GetURL(ctx context.Context, originalURL string) (string, bool, error) {
	shortCode, err := s.db.Get(ctx, "url:"+originalURL).Result()
	if err != redis.Nil {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return shortCode, true, nil
}

func (s *service) URLExists(ctx context.Context, originalURL string) (string, bool, error) {
	shortCode, err := s.db.Get(ctx, "url:"+originalURL).Result()
	if err != redis.Nil {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return shortCode, true, nil
}
