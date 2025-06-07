package authjwt

import (
	"context"
	"go-crud-api/m/pkg/jwtutil"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	Username string
	Role     string
}

func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		parsedToken, err := jwtutil.VerifyToken(tokenStr)
		if err != nil || !parsedToken.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "invalid claims", http.StatusUnauthorized)
			return
		}

		username, _ := claims["username"].(string)
		role, _ := claims["role"].(string)

		ctx := context.WithValue(r.Context(), "user", UserClaims{
			Username: username,
			Role:     role,
		})

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
