package middleware

import (
	"context"
	"net/http"
	"strings"

	"digital-greenhouse/greenhouse-be/internal/security"
)

type contextKey string

const (
	UserIDKey contextKey = "userID"
	UserRoleKey contextKey = "userRole"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Se requiere token de autenticación", http.StatusUnauthorized)
			return
		}

		// El formato esperado es "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Formato de token inválido", http.StatusUnauthorized)
			return
		}

		claims, err := security.ValidateToken(parts[1])
		if err != nil {
			http.Error(w, "Token inválido o expirado: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// Extraer claims (sub es el userID en nuestro security/jwt.go)
		userIDFloat, ok := claims["sub"].(float64)
		if !ok {
			http.Error(w, "Token no contiene ID de usuario", http.StatusUnauthorized)
			return
		}

		role, _ := claims["role"].(string)

		// Inyectar en el contexto
		ctx := context.WithValue(r.Context(), UserIDKey, uint(userIDFloat))
		ctx = context.WithValue(ctx, UserRoleKey, role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func OptionalAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			next.ServeHTTP(w, r)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			// Si envió algo pero está mal formateado, lo tratamos como no autenticado
			// Opcionalmente podríamos fallar, pero por simplicidad para el front,
			// si no es un Bearer token válido, no inyectamos datos.
			next.ServeHTTP(w, r)
			return
		}

		claims, err := security.ValidateToken(parts[1])
		if err != nil {
			// Token inválido (expirado, malformado), simplemente no inyectamos usuario
			next.ServeHTTP(w, r)
			return
		}

		userIDFloat, ok := claims["sub"].(float64)
		if !ok {
			next.ServeHTTP(w, r)
			return
		}

		role, _ := claims["role"].(string)

		ctx := context.WithValue(r.Context(), UserIDKey, uint(userIDFloat))
		ctx = context.WithValue(ctx, UserRoleKey, role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Helpers para obtener datos del contexto en los handlers
func GetUserID(ctx context.Context) uint {
	val, _ := ctx.Value(UserIDKey).(uint)
	return val
}

func GetUserRole(ctx context.Context) string {
	val, _ := ctx.Value(UserRoleKey).(string)
	return val
}
