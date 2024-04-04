package middlewares

import (
	"log"
	"net/http"
	"strings"
)

func ConnectionLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr[:strings.LastIndex(r.RemoteAddr, ":")]
		log.Printf("[ %s ] %s %s\n", ip, r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
