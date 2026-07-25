package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/yourorg/go-user-service/internal/domain/session"
)

const sessionKeyPrefix = "auth:session"

// ErrSessionNotFound is returned when a session ID has no matching data.
var ErrSessionNotFound = errors.New("session not found")

type sessionStore struct {
	client *redis.Client
	ttl    time.Duration
}

// NewSessionStore wires a Redis client into the session.Store port.
func NewSessionStore(client *redis.Client, ttl time.Duration) session.Store {
	return &sessionStore{client: client, ttl: ttl}
}

func (s *sessionStore) buildKey(sessionID string) string {
	return fmt.Sprintf("%s:%s", sessionKeyPrefix, sessionID)
}

func (s *sessionStore) Set(ctx context.Context, sessionID string, data session.Data) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal session data: %w", err)
	}
	if err := s.client.Set(ctx, s.buildKey(sessionID), payload, s.ttl).Err(); err != nil {
		return fmt.Errorf("failed to set session in redis: %w", err)
	}
	return nil
}

func (s *sessionStore) Get(ctx context.Context, sessionID string) (*session.Data, error) {
	payload, err := s.client.Get(ctx, s.buildKey(sessionID)).Result()
	if errors.Is(err, redis.Nil) {
		return nil, ErrSessionNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get session from redis: %w", err)
	}

	var data session.Data
	if err := json.Unmarshal([]byte(payload), &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session data: %w", err)
	}
	return &data, nil
}

func (s *sessionStore) Delete(ctx context.Context, sessionID string) error {
	if err := s.client.Del(ctx, s.buildKey(sessionID)).Err(); err != nil {
		return fmt.Errorf("failed to delete session from redis: %w", err)
	}
	return nil
}