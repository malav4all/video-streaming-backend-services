package session

import "context"

// StateStore manages the short-lived OAuth2 "state" parameter used to
// prevent CSRF attacks during the authorization code flow.
type StateStore interface {
	SetState(ctx context.Context, state string) error
	GetState(ctx context.Context, state string) (string, error)
	DeleteState(ctx context.Context, state string) error
}

// Store manages authenticated user sessions, keyed by session ID.
type Store interface {
	Set(ctx context.Context, sessionID string, data Data) error
	Get(ctx context.Context, sessionID string) (*Data, error)
	Delete(ctx context.Context, sessionID string) error
}