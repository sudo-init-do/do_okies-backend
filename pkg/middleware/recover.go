package middleware

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Error().Interface("panic", rec).Msg("panic recovered")
				http.Error(w, "internal_error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
