package random

import (
	"crypto/rand"
	"encoding/base64"
)

// SecureToken returns a cryptographically secure, URL-safe random string.
// Used for OAuth2 state parameters and session IDs.
func SecureToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(buf), nil
}