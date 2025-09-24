package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	UserContextKey contextKey = "user"
)

type UserContext struct {
	UserID   string `json:"userId"`
	Role     string `json:"role"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	FullName string `json:"fullName"`
	UserType string `json:"userType"`
}

type JWTMiddleware struct {
	secretKey []byte
}

func NewJWTMiddleware(secretKey string) *JWTMiddleware {
	return &JWTMiddleware{
		secretKey: []byte(secretKey),
	}
}

func (m *JWTMiddleware) ValidateToken(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := extractTokenFromHeader(r)
		if tokenString == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		userCtx, err := m.parseToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Set headers for handlers to use
		r.Header.Set("X-User-ID", userCtx.UserID)
		r.Header.Set("X-User-Email", userCtx.Email)
		r.Header.Set("X-User-Role", userCtx.Role)
		r.Header.Set("X-User-Name", userCtx.Name)

		ctx := context.WithValue(r.Context(), UserContextKey, userCtx)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func (m *JWTMiddleware) OptionalAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := extractTokenFromHeader(r)
		if tokenString != "" {
			if userCtx, err := m.parseToken(tokenString); err == nil {
				// Set headers for handlers to use
				r.Header.Set("X-User-ID", userCtx.UserID)
				r.Header.Set("X-User-Email", userCtx.Email)
				r.Header.Set("X-User-Role", userCtx.Role)
				r.Header.Set("X-User-Name", userCtx.Name)

				ctx := context.WithValue(r.Context(), UserContextKey, userCtx)
				r = r.WithContext(ctx)
			}
		}
		next.ServeHTTP(w, r)
	}
}

func (m *JWTMiddleware) parseToken(tokenString string) (*UserContext, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	userCtx := &UserContext{}
	
	// Extract user data from nested "user" object in JWT
	if userObj, ok := claims["user"].(map[string]interface{}); ok {
		if userID, ok := userObj["id"].(string); ok {
			userCtx.UserID = userID
		}
		if email, ok := userObj["email"].(string); ok {
			userCtx.Email = email
		}
		if name, ok := userObj["name"].(string); ok {
			userCtx.Name = name
		}
		if fullName, ok := userObj["fullName"].(string); ok {
			userCtx.FullName = fullName
		}
		if userType, ok := userObj["userType"].(string); ok {
			userCtx.UserType = userType
			// Map userType to role for backwards compatibility
			switch userType {
			case "admin":
				userCtx.Role = "admin"
			case "instructor":
				userCtx.Role = "instructor"
			default:
				userCtx.Role = "student"
			}
		}
	} else {
		// Fallback: try to extract from top-level claims
		if userID, ok := claims["userId"].(string); ok {
			userCtx.UserID = userID
		} else if sub, ok := claims["sub"].(string); ok {
			userCtx.UserID = sub
		}
		
		if role, ok := claims["role"].(string); ok {
			userCtx.Role = role
		}
		
		if email, ok := claims["email"].(string); ok {
			userCtx.Email = email
		}
	}

	if userCtx.UserID == "" {
		return nil, fmt.Errorf("user ID not found in token")
	}

	return userCtx, nil
}

func extractTokenFromHeader(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}

	return parts[1]
}

func GetUserFromContext(ctx context.Context) (*UserContext, bool) {
	user, ok := ctx.Value(UserContextKey).(*UserContext)
	return user, ok
}

func RequireRole(allowedRoles ...string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			user, ok := GetUserFromContext(r.Context())
			if !ok {
				http.Error(w, "User context not found", http.StatusUnauthorized)
				return
			}

			roleAllowed := false
			for _, role := range allowedRoles {
				if user.Role == role {
					roleAllowed = true
					break
				}
			}

			if !roleAllowed {
				http.Error(w, "Insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		}
	}
}