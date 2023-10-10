package handler

import (
	"net/http"
)

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers to allow all origins
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Handle preflight requests (OPTIONS)
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			return
		}

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}
