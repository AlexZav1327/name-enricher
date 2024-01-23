package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (h *Handler) metric(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		started := time.Now()

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)
		pattern := chi.RouteContext(r.Context()).RoutePattern()

		h.metrics.requests.WithLabelValues(http.StatusText(ww.Status()), r.Method, pattern).Inc()
		h.metrics.duration.WithLabelValues(http.StatusText(ww.Status()), r.Method,
			pattern).Observe(time.Since(started).Seconds())
	}

	return http.HandlerFunc(fn)
}

func (*Handler) commonMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
