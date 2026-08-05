package session

import "time"

// UserInfo holds the identity claims extracted from the Keycloak ID token.
type UserInfo struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

// Data represents an authenticated user's session, persisted in the
// session store and referenced by a client-side session cookie.
type Data struct {
	AccessToken string    `json:"access_token"`
	UserInfo    UserInfo  `json:"user_info"`
	CreatedAt   time.Time `json:"created_at"`
}