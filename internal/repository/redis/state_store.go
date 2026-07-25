package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/yourorg/go-user-service/internal/domain/session"
)

const (
	stateKeyPrefix = "auth:state"
	stateTTL       = 2 * time.Minute
)

type stateStore struct {
	client *redis.Client
}

// NewStateStore wires a Redis client into the session.StateStore port.
func NewStateStore(client *redis.Client) session.StateStore {
	return &stateStore{client: client}
}

func (s *stateStore) buildKey(state string) string {
	return fmt.Sprintf("%s:%s", stateKeyPrefix, state)
}

func (s *stateStore) SetState(ctx context.Context, state string) error {
	if err := s.client.Set(ctx, s.buildKey(state), state, stateTTL).Err(); err != nil {
		return fmt.Errorf("failed to set state in redis: %w", err)
	}
	return nil
}

func (s *stateStore) GetState(ctx context.Context, state string) (string, error) {
	value, err := s.client.Get(ctx, s.buildKey(state)).Result()
	if err != nil {
		return "", fmt.Errorf("failed to get state from redis: %w", err)
	}
	return value, nil
}

func (s *stateStore) DeleteState(ctx context.Context, state string) error {
	if err := s.client.Del(ctx, s.buildKey(state)).Err(); err != nil {
		return fmt.Errorf("failed to delete state from redis: %w", err)
	}
	return nil
}