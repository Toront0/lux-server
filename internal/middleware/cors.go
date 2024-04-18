package middleware

import (
	"net/http"
)

type Middleware func (http.Handler) http.Handler

func CreateStack(xs ...Middleware) Middleware {
	return func (next http.Handler) http.Handler {

		for i := len(xs) - 1; i >= 0; i-- {

			x := xs[i]
			next = x(next)
		}

		return next
	}
}

func AllowCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "http://*")
		
		// w.Header().Set("Access-Control-Max-Age", "300")
		w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Accept, Content-Type, X-CSRF-Token, Authorization",)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Expose-Headers", "Link")
		// w.Header().Set("Allowed")
		// AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		next.ServeHTTP(w, r)

	})


}