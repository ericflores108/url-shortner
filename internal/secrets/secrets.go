package secrets

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type RedisConfig struct {
	Address  string `json:"address"`
	Port     string `json:"port"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type SecretsManager struct {
	client *secretsmanager.Client
}

func NewSecretsManager(ctx context.Context) (*SecretsManager, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := secretsmanager.NewFromConfig(cfg)
	return &SecretsManager{client: client}, nil
}

func (sm *SecretsManager) GetRedisConfig(ctx context.Context, secretName string) (*RedisConfig, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId: &secretName,
	}

	result, err := sm.client.GetSecretValue(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get secret %s: %w", secretName, err)
	}

	var redisConfig RedisConfig
	if err := json.Unmarshal([]byte(*result.SecretString), &redisConfig); err != nil {
		return nil, fmt.Errorf("failed to unmarshal secret: %w", err)
	}

	return &redisConfig, nil
}

func GetRedisConfigForEnv() (*RedisConfig, error) {
	// Check if we're in production (Lambda environment)
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" || os.Getenv("ENV") == "prod" {
		ctx := context.Background()
		sm, err := NewSecretsManager(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to create secrets manager: %w", err)
		}

		secretName := "prod/urlshort/redis"
		config, err := sm.GetRedisConfig(ctx, secretName)
		if err != nil {
			return nil, fmt.Errorf("failed to get Redis config from secrets manager: %w", err)
		}

		return config, nil
	}

	// For local development, use environment variables
	return &RedisConfig{
		Address:  os.Getenv("URLSHORT_DB_ADDRESS"),
		Port:     os.Getenv("URLSHORT_DB_PORT"),
		Password: os.Getenv("URLSHORT_DB_PASSWORD"),
		Database: os.Getenv("URLSHORT_DB_DATABASE"),
	}, nil
}
