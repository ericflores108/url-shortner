package database

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/redis"
)

func mustStartRedisContainer() (func(context.Context, ...testcontainers.TerminateOption) error, error) {
	dbContainer, err := redis.Run(
		context.Background(),
		"docker.io/redis:8.0.2",
		redis.WithSnapshotting(10, 1),
		redis.WithLogLevel(redis.LogLevelVerbose),
	)
	if err != nil {
		return nil, err
	}

	dbHost, err := dbContainer.Host(context.Background())
	if err != nil {
		return dbContainer.Terminate, err
	}

	dbPort, err := dbContainer.MappedPort(context.Background(), "6379/tcp")
	if err != nil {
		return dbContainer.Terminate, err
	}

	os.Setenv("URLSHORT_DB_ADDRESS", dbHost)
	os.Setenv("URLSHORT_DB_PORT", dbPort.Port())
	os.Setenv("URLSHORT_DB_DATABASE", "0")
	os.Setenv("URLSHORT_DB_PASSWORD", "")

	return dbContainer.Terminate, err
}

func TestMain(m *testing.M) {
	teardown, err := mustStartRedisContainer()
	if err != nil {
		log.Fatalf("could not start redis container: %v", err)
	}

	m.Run()

	if teardown != nil && teardown(context.Background()) != nil {
		log.Fatalf("could not teardown redis container: %v", err)
	}
}

func TestNew(t *testing.T) {
	srv := New()
	if srv == nil {
		t.Fatal("New() returned nil")
	}
}

func TestHealth(t *testing.T) {
	srv := New()

	stats := srv.Health()

	if stats["redis_status"] != "up" {
		t.Fatalf("expected status to be up, got %s", stats["redis_status"])
	}

	if _, ok := stats["redis_version"]; !ok {
		t.Fatalf("expected redis_version to be present, got %v", stats["redis_version"])
	}
}
