package middleware

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type lrw struct {
	http.ResponseWriter
	status int
}

func (w *lrw) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func RequestLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		l := &lrw{ResponseWriter: w, status: 200}
		next.ServeHTTP(l, r)
		dur := time.Since(start)
		ev := log.Info()
		if l.status >= 400 && l.status < 500 { ev = log.Warn() }
		if l.status >= 500 { ev = log.Error() }
		ev.Str("method", r.Method).
			Str("url", r.URL.Path).
			Int("status", l.status).
			Float64("duration_ms", float64(dur.Microseconds())/1000.0).
			Msg("request completed")
	})
}
