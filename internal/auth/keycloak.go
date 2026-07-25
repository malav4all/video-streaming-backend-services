package auth

import (
	"context"
	"fmt"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

// Config holds the Keycloak/OIDC client settings needed to talk to a realm.
type Config struct {
	BaseURL      string
	Realm        string
	ClientID     string
	ClientSecret string					
	RedirectURL  string
}

// Client bundles everything needed to run the OAuth2/OIDC authorization
// code flow against Keycloak: the OIDC provider, an ID token verifier,
// and the OAuth2 client configuration.
type Client struct {
	Provider *oidc.Provider
	Verifier *oidc.IDTokenVerifier
	OAuth2   oauth2.Config
}

// NewKeycloakClient discovers the given realm's OIDC provider and builds
// an OAuth2 client configured for the authorization code flow.
func NewKeycloakClient(ctx context.Context, cfg *Config) (*Client, error) {
	providerURL := fmt.Sprintf("%s/realms/%s", cfg.BaseURL, cfg.Realm)

	provider, err := oidc.NewProvider(ctx, providerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to discover keycloak provider: %w", err)
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: cfg.ClientID})

	oauth2Config := oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email", "roles"},
	}

	return &Client{
		Provider: provider,
		Verifier: verifier,
		OAuth2:   oauth2Config,
	}, nil
}