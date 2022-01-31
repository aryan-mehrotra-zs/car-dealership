package middlewares

import (
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header.Get("Api-Key") != "aryan-zs" {
			w.WriteHeader(http.StatusUnauthorized)

			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)

	})
}
