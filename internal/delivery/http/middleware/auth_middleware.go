package middleware

import (
	"errors"
	"fmt"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gofiber/fiber/v2"

	"github.com/yourorg/go-user-service/internal/auth"
	"github.com/yourorg/go-user-service/internal/delivery/http/response"
	"github.com/yourorg/go-user-service/internal/domain/session"
	"github.com/yourorg/go-user-service/internal/domain/tenant"
)

const (
	sessionCookieName = "session_id"
	LocalsSessionKey  = "user_session"
	LocalsClaimsKey   = "user_claims"
	LocalsTenantIDKey = "tenant_id"
	LocalsRolesKey    = "user_roles"
)

// AuthMiddleware guards routes behind an authenticated Keycloak session.
type AuthMiddleware struct {
	authClient   *auth.Client
	sessionStore session.Store
	tenantSvc    tenant.Service
}

func NewAuthMiddleware(authClient *auth.Client, sessionStore session.Store, tenantSvc tenant.Service) *AuthMiddleware {
	return &AuthMiddleware{authClient: authClient, sessionStore: sessionStore, tenantSvc: tenantSvc}
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

		tenantID, err := m.resolveTenantID(c, claims)
		if err != nil {
			return response.Error(c, fiber.StatusForbidden, "no tenant association found", nil)
		}

		roles := m.extractRoles(claims)

		c.Locals(LocalsSessionKey, sessionData)
		c.Locals(LocalsClaimsKey, claims)
		c.Locals(LocalsTenantIDKey, tenantID)
		c.Locals(LocalsRolesKey, roles)

		return c.Next()
	}
}

// RequireRole guards a route to only the listed realm roles. Must run after
// RequireAuth(), since it reads roles from c.Locals set there.
func (m *AuthMiddleware) RequireRole(allowed ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		roles, ok := c.Locals(LocalsRolesKey).([]string)
		if !ok {
			return response.Error(c, fiber.StatusForbidden, "no roles found on session", nil)
		}

		for _, role := range roles {
			for _, want := range allowed {
				if role == want {
					return c.Next()
				}
			}
		}

		return response.Error(c, fiber.StatusForbidden, "insufficient permissions", nil)
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

// resolveTenantID reads the first Keycloak group from claims (e.g. "/tenant-acme"),
// strips the leading slash, and resolves it to a real tenant UUID.
func (m *AuthMiddleware) resolveTenantID(c *fiber.Ctx, claims map[string]interface{}) (string, error) {
	groupsRaw, ok := claims["groups"].([]interface{})
	if !ok || len(groupsRaw) == 0 {
		return "", errors.New("no groups claim present")
	}

	groupPath, ok := groupsRaw[0].(string)
	if !ok {
		return "", errors.New("invalid groups claim format")
	}

	slug := strings.TrimPrefix(groupPath, "/")

	t, err := m.tenantSvc.GetBySlug(c.Context(), slug)
	if err != nil {
		return "", fmt.Errorf("failed to resolve tenant for slug %q: %w", slug, err)
	}

	return t.ID, nil
}

// extractRoles pulls realm_access.roles out of the claims map.
func (m *AuthMiddleware) extractRoles(claims map[string]interface{}) []string {
	realmAccess, ok := claims["realm_access"].(map[string]interface{})
	if !ok {
		return nil
	}
	rolesRaw, ok := realmAccess["roles"].([]interface{})
	if !ok {
		return nil
	}

	roles := make([]string, 0, len(rolesRaw))
	for _, r := range rolesRaw {
		if s, ok := r.(string); ok {
			roles = append(roles, s)
		}
	}
	return roles
}
