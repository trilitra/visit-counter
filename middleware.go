package main

import (
	"context"
	"net/http"
)

func ipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ip, err := clientIP(r)
		if err != nil {
			http.Error(w, "bad remote addr", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), ipKey, ip)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
