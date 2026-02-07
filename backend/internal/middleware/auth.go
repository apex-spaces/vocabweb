package middleware

import (
	"context"
	"net/http"
	"strings"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
)

type contextKey string

const (
	UserIDKey contextKey = "userID"
	EmailKey  contextKey = "email"
)

type AuthMiddleware struct {
	authClient *auth.Client
}

func NewAuthMiddleware(ctx context.Context, projectID string) (*AuthMiddleware, error) {
	app, err := firebase.NewApp(ctx, &firebase.Config{
		ProjectID: projectID,
	})
	if err != nil {
		return nil, err
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		return nil, err
	}

	return &AuthMiddleware{authClient: authClient}, nil
}

func (am *AuthMiddleware) Verify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}

		verifiedToken, err := am.authClient.VerifyIDToken(r.Context(), token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Extract email from token claims
		email := ""
		if emailClaim, ok := verifiedToken.Claims["email"].(string); ok {
			email = emailClaim
		}

		// Set user_id and email to context
		ctx := context.WithValue(r.Context(), UserIDKey, verifiedToken.UID)
		ctx = context.WithValue(ctx, EmailKey, email)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
