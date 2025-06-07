package database

import (
	"context"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

func (s *service) StoreURL(ctx context.Context, shortCode, originalURL string) error {
	const oneWeek = 7 * 24 * time.Hour

	err := s.db.Set(ctx, "code:"+shortCode, originalURL, oneWeek).Err()
	if err != nil {
		return err
	}

	err = s.db.Set(ctx, "url:"+originalURL, shortCode, oneWeek).Err()
	if err != nil {
		return err
	}

	return nil
}

func (s *service) GetURL(ctx context.Context, originalURL string) (string, bool, error) {
	shortCode, err := s.db.Get(ctx, "url:"+originalURL).Result()
	switch {
	case err == redis.Nil:
		return "", false, nil
	case err != nil:
		slog.Warn("Get failed", "err", err)
		return "", false, err
	case shortCode == "":
		slog.Debug("shortCode is empty", "originalURL", shortCode)
		return "", false, nil
	}
	return shortCode, true, nil
}

func (s *service) URLExists(ctx context.Context, originalURL string) (string, bool, error) {
	shortCode, err := s.db.Get(ctx, "url:"+originalURL).Result()
	if err == redis.Nil {
		return "", false, nil
	}
	if err != nil {
		slog.Warn("Get failed", "err", err)
		return "", false, err
	}
	return shortCode, true, nil
}
