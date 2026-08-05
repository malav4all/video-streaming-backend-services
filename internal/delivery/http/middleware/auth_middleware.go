package middleware

import (
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gofiber/fiber/v2"

	"github.com/yourorg/go-user-service/internal/auth"
	"github.com/yourorg/go-user-service/internal/delivery/http/response"
	"github.com/yourorg/go-user-service/internal/domain/session"
)

const (
	sessionCookieName = "session_id"
	LocalsSessionKey  = "user_session"
	LocalsClaimsKey   = "user_claims"
)

// AuthMiddleware guards routes behind an authenticated Keycloak session.
type AuthMiddleware struct {
	authClient   *auth.Client
	sessionStore session.Store
}

func NewAuthMiddleware(authClient *auth.Client, sessionStore session.Store) *AuthMiddleware {
	return &AuthMiddleware{authClient: authClient, sessionStore: sessionStore}
}

// RequireAuth validates the session cookie, confirms the underlying access
// token is still valid against Keycloak, and injects the session and token
// claims into the request context (via c.Locals) for downstream handlers.
func (m *AuthMiddleware) RequireAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		sessionID := c.Cookies(sessionCookieName)
		if sessionID == "" {
			return response.Error(c, fiber.StatusUnauthorized, "authentication required", nil)
		}

		sessionData, err := m.sessionStore.Get(c.Context(), sessionID)
		if err != nil {
			m.clearSessionCookie(c)
			return response.Error(c, fiber.StatusUnauthorized, "invalid or expired session", nil)
		}

		verifier := m.authClient.Provider.Verifier(&oidc.Config{SkipClientIDCheck: true})
		token, err := verifier.Verify(c.Context(), sessionData.AccessToken)
		if err != nil {
			_ = m.sessionStore.Delete(c.Context(), sessionID)
			m.clearSessionCookie(c)
			return response.Error(c, fiber.StatusUnauthorized, "access token is no longer valid", nil)
		}

		var claims map[string]interface{}
		if err := token.Claims(&claims); err != nil {
			return response.Error(c, fiber.StatusInternalServerError, "failed to parse token claims", nil)
		}

		c.Locals(LocalsSessionKey, sessionData)
		c.Locals(LocalsClaimsKey, claims)

		return c.Next()
	}
}

func (m *AuthMiddleware) clearSessionCookie(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HTTPOnly: true,
		Secure:   true,
		SameSite: fiber.CookieSameSiteStrictMode,
	})
}