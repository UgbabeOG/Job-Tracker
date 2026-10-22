package httpapi

import (
	"net/http"
	"strings"
)

func isAllowedMethod(method string) bool {
	switch method {
	case http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodDelete:
		return true
	default:
		return false
	}
}

func CORSMiddleware(next http.Handler, allowedOrigins []string) http.Handler {
	allowed := make(map[string]bool, len(allowedOrigins))

	for _, origin := range allowedOrigins {
		origin = strings.TrimSpace(origin)
		if origin != "" {
			allowed[origin] = true
		}
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if origin == "" {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Add("Vary", "Origin")

		if !allowed[origin] {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", origin)

		if r.Method == http.MethodOptions {
			requestedMethod := r.Header.Get("Access-Control-Request-Method")
			if !isAllowedMethod(requestedMethod) {
				http.Error(w, "CORS method not allowed", http.StatusForbidden)
				return
			}
			w.Header().Set(
				"Access-Control-Allow-Methods",
				"GET, POST, PATCH, DELETE, OPTIONS",
			)
			w.Header().Set(
				"Access-Control-Allow-Headers",
				"Content-Type, Authorization",
			)
			w.Header().Set("Access-Control-Max-Age", "600")
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
