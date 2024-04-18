package middleware

import (
	"net/http"
	"fmt"
	"context"
	"slices"

	"guthub.com/Toront0/lux-server/internal/utils"
	"github.com/golang-jwt/jwt/v5"


)

type AuthContextUserID int

const CtxAuthKey AuthContextUserID = 0

func UserFromContext(ctx context.Context) int {

	val := ctx.Value(CtxAuthKey).(float64)

	return int(val)
}

var excludeAuthPaths = []string{"/auth", "/login", "/sign-up", "/users/7/music-playlists"}

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		skip := slices.Contains(excludeAuthPaths, r.URL.Path)


		if skip {
			next.ServeHTTP(w, r)
			return
		}
		
		

		cookie, err := r.Cookie("jwt")
	
	
	

		if err != nil {
			w.WriteHeader(401)
			return
		}

		token, err := utils.ValidateJWT(cookie.Value)

		if err != nil {
			fmt.Printf("invalid JWT Token %s", err)
			w.WriteHeader(400)
			return
		}

		

		claims := token.Claims.(jwt.MapClaims)

		userID := claims["userID"]


		ctx := context.WithValue(r.Context(), CtxAuthKey, userID)
		
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}