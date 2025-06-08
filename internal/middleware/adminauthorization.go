package middleware

import "net/http"

func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		claims, ok := request.Context().Value("user").(UserClaims)
		if !ok || claims.Role != "admin" {
			http.Error(response, "access denied: admin only", http.StatusForbidden)
			return
		}
		next.ServeHTTP(response, request)
	})
}
