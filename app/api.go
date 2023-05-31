package app

import (
	"io"
	"net/http"
	"port-logger/pkg/logger"
)

func WithLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		userAgent := r.UserAgent()
		logger.Infof(ctx, "[WithLogger] Response from %s", userAgent)

		for name, values := range r.Header {
			// Loop over all values for the name.
			for _, value := range values {
				logger.Infof(ctx, "[WithLogger] Header| %s: %s", name, value)
			}
		}

		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Errorf(ctx, "[WithLogger] Response body read error: %s", err)
		} else {
			bodyString := string(bodyBytes)
			logger.Infof(ctx, "[WithLogger] Response body: %s", bodyString)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Server) Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}
