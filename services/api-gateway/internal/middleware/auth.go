package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	pb "github.com/maxluvds/booking-platform/api/proto"
	"github.com/maxluvds/booking-platform/services/api-gateway/internal/client"
)

type contextKey string

const UserIDKey contextKey = "user_id"
const UserEmailKey contextKey = "user_email"
const UserRoleKey contextKey = "user_role"

func JWTAuth(userClient *client.UserClient, log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				log.Warn("Missing authorization header")
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				log.Warn("Invalid authorization header format", "header", authHeader)
				http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
				return
			}

			token := parts[1]
			log.Info("Validating token", "token_prefix", token[:20])

			resp, err := userClient.ValidateToken(r.Context(), &pb.ValidateTokenRequest{Token: token})
			if err != nil {
				log.Error("Failed to validate token", "error", err)
				http.Error(w, "Failed to validate token", http.StatusUnauthorized)
				return
			}

			log.Info("Token validation result", "valid", resp.Valid)

			if !resp.Valid {
				log.Warn("Token is invalid")
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, resp.UserId)
			ctx = context.WithValue(ctx, UserEmailKey, resp.Email)
			ctx = context.WithValue(ctx, UserRoleKey, resp.Role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
