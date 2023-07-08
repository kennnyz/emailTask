package delivery

import (
	"log"
	"net/http"
)

func staticAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//authHeader := r.Header.Get("Authorization")
		//if authHeader != "secret" {
		//	http.Error(w, "unauthorized", http.StatusUnauthorized)
		//	return
		//}
		log.Println("auth success ")
		next.ServeHTTP(w, r)
	})
}
