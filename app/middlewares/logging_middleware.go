package middlewares

import (
	"log"
	"net/http"
	"time"
)

type responseWriterWithCode struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriterWithCode) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func LogAccessMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rw := &responseWriterWithCode{ResponseWriter: w, statusCode: http.StatusOK}
			next.ServeHTTP(rw, r)

			clientIP := r.RemoteAddr
			if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
				clientIP = forwarded
			}
			log.Printf(
				"%s %s %s %d %s %s %s",
				clientIP, r.Method, r.URL.Path, rw.statusCode, time.Since(start), r.UserAgent(), r.URL.RawQuery,
			)
		},
	)
}

func RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					w.Header().Set("Connection", "close")

					log.Printf("Internal Server Error, Application recovered, eror: %s", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		},
	)
}
