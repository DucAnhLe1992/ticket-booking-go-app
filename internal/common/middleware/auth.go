package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const userContextKey contextKey = "currentUser"

// UserClaims captures the JWT claims we care about.
type UserClaims struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// CurrentUser extracts a JWT from headers/cookies and attaches claims to context.
func CurrentUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := extractToken(r)
		if tokenStr == "" {
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userContextKey, (*UserClaims)(nil))))
			return
		}

		key := os.Getenv("JWT_KEY")
		if key == "" {
			// If JWT key missing, treat as unauthenticated environment (dev).
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userContextKey, (*UserClaims)(nil))))
			return
		}

		token, _ := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(key), nil
		})

		var claims *UserClaims
		if token != nil {
			if c, ok := token.Claims.(*UserClaims); ok && token.Valid {
				claims = c
			}
		}

		ctx := context.WithValue(r.Context(), userContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireAuth enforces presence of a valid current user.
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if cu := GetCurrentUser(r.Context()); cu == nil || cu.ID == "" {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"message":"Not authorized"}`))
			return
		}
		next.ServeHTTP(w, r)
	})
}

// GetCurrentUser retrieves claims from context.
func GetCurrentUser(ctx context.Context) *UserClaims {
	if v := ctx.Value(userContextKey); v != nil {
		if c, ok := v.(*UserClaims); ok {
			return c
		}
	}
	return nil
}

// extractToken tries common locations for JWT: Authorization header, cookie "jwt",
// or cookie "session" containing a JSON object {"jwt":"..."} as used in the TS app.
func extractToken(r *http.Request) string {
	// 1) Authorization: Bearer <token>
	if h := r.Header.Get("Authorization"); h != "" {
		parts := strings.SplitN(h, " ", 2)
		if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
			return parts[1]
		}
	}
	// 2) Cookie: jwt
	if c, err := r.Cookie("jwt"); err == nil {
		if c.Value != "" {
			return c.Value
		}
	}
	// 3) Cookie: session with JSON containing jwt field
	if c, err := r.Cookie("session"); err == nil && c.Value != "" {
		var obj struct {
			JWT string `json:"jwt"`
		}
		_ = json.Unmarshal([]byte(c.Value), &obj)
		if obj.JWT != "" {
			return obj.JWT
		}
	}
	return ""
}
