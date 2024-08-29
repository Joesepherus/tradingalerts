package authMiddleware

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func TokenAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenString := cookie.Value
		claims := jwt.MapClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("your-secret-key"), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			exp, ok := claims["exp"].(float64)
			if !ok || time.Now().Unix() > int64(exp) {
				http.Error(w, "Token expired", http.StatusUnauthorized)
				return
			}
		}

		// Token is valid, pass the request to the next handler
		next.ServeHTTP(w, r)
	})
}
