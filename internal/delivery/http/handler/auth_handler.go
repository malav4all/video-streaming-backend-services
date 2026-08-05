package handler

import (
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"

	"github.com/yourorg/go-user-service/internal/auth"
	"github.com/yourorg/go-user-service/internal/delivery/http/response"
	"github.com/yourorg/go-user-service/internal/domain/session"
	"github.com/yourorg/go-user-service/pkg/random"
)

const sessionCookieName = "session_id"

// AuthHandler drives the Keycloak OAuth2/OIDC authorization code flow:
// login initiation, callback processing, and logout.
type AuthHandler struct {
	authClient   *auth.Client
	stateStore   session.StateStore
	sessionStore session.Store
	sessionTTL   time.Duration
}

func NewAuthHandler(
	authClient *auth.Client,
	stateStore session.StateStore,
	sessionStore session.Store,
	sessionTTL time.Duration,
) *AuthHandler {
	return &AuthHandler{
		authClient:   authClient,
		stateStore:   stateStore,
		sessionStore: sessionStore,
		sessionTTL:   sessionTTL,
	}
}

// Login redirects the user to Keycloak's login page, having first stored
// a CSRF-protection state value in the state store.
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	state, err := random.SecureToken()
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to generate state", nil)
	}

	if err := h.stateStore.SetState(c.Context(), state); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to persist state", nil)
	}

	authURL := h.authClient.OAuth2.AuthCodeURL(state, oauth2.SetAuthURLParam("response_type", "code"))

	return c.Redirect(authURL, fiber.StatusTemporaryRedirect)
}

// Callback handles Keycloak's redirect after successful authentication:
// validates state, exchanges the code for tokens, verifies the ID token,
// and establishes a server-side session referenced by a cookie.
func (h *AuthHandler) Callback(c *fiber.Ctx) error {
	if err := h.validateState(c); err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "invalid state parameter", nil)
	}

	oauthToken, err := h.exchangeCode(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "failed to exchange authorization code", nil)
	}

	claims, err := h.verifyIDToken(c, oauthToken)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "failed to verify id token", nil)
	}

	sessionID, err := random.SecureToken()
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to generate session id", nil)
	}

	sessionData := session.Data{
		AccessToken: oauthToken.AccessToken,
		UserInfo: session.UserInfo{
			Username: claims.Username,
			Email:    claims.Email,
		},
		CreatedAt: time.Now(),
	}

	if err := h.sessionStore.Set(c.Context(), sessionID, sessionData); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to persist session", nil)
	}

	c.Cookie(&fiber.Cookie{
		Name:     sessionCookieName,
		Value:    sessionID,
		Path:     "/",
		MaxAge:   int(h.sessionTTL.Seconds()),
		HTTPOnly: true,
		Secure:   true,
		SameSite: fiber.CookieSameSiteStrictMode,
	})

	return response.Success(c, fiber.StatusOK, "login successful", sessionData.UserInfo)
}

// Logout clears the session both server-side and via the client cookie.
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	sessionID := c.Cookies(sessionCookieName)
	if sessionID != "" {
		_ = h.sessionStore.Delete(c.Context(), sessionID)
	}

	h.clearSessionCookie(c)

	return response.Success(c, fiber.StatusOK, "logout successful", nil)
}

func (h *AuthHandler) validateState(c *fiber.Ctx) error {
	state := c.Query("state")
	if state == "" {
		return errors.New("missing state parameter")
	}

	storedState, err := h.stateStore.GetState(c.Context(), state)
	if err != nil {
		return fmt.Errorf("failed to retrieve stored state: %w", err)
	}

	if storedState != state {
		return errors.New("state parameter mismatch")
	}

	_ = h.stateStore.DeleteState(c.Context(), state)
	return nil
}

func (h *AuthHandler) exchangeCode(c *fiber.Ctx) (*oauth2.Token, error) {
	code := c.Query("code")
	if code == "" {
		return nil, errors.New("missing authorization code")
	}

	return h.authClient.OAuth2.Exchange(
		c.Context(), code,
		oauth2.SetAuthURLParam("grant_type", "authorization_code"),
	)
}

// idTokenClaims are the subset of Keycloak ID token claims we care about.
type idTokenClaims struct {
	Email    string `json:"email"`
	Username string `json:"preferred_username"`
}

func (h *AuthHandler) verifyIDToken(c *fiber.Ctx, token *oauth2.Token) (*idTokenClaims, error) {
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("no id_token found in token response")
	}

	idToken, err := h.authClient.Verifier.Verify(c.Context(), rawIDToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify id token: %w", err)
	}

	var claims idTokenClaims
	if err := idToken.Claims(&claims); err != nil {
		return nil, fmt.Errorf("failed to parse id token claims: %w", err)
	}

	return &claims, nil
}

func (h *AuthHandler) clearSessionCookie(c *fiber.Ctx) {
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