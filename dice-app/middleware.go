package main

import (
	"net/http"
	"strings"
	"time"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/attribute"
)

type Middleware func(http.Handler) http.Handler

var publicMiddlewares = []Middleware{
	TelemetryMiddleware,
}

var protectedMiddlewares = []Middleware{
	RequireAuth,
	TelemetryMiddleware,
}

func Chain(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}

func getSessionID(r *http.Request) string {
	// Authorization: Bearer <sessionID>
	if h := r.Header.Get("Authorization"); h != "" {
		if strings.HasPrefix(h, "Bearer ") {
			return strings.TrimSpace(h[7:])
		}
	}
	// X-Session-ID header
	if h := r.Header.Get("X-Session-ID"); h != "" {
		return strings.TrimSpace(h)
	}
	// Cookie session_id
	if c, err := r.Cookie("session_id"); err == nil && c.Value != "" {
		return c.Value
	}
	return ""
}

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authLogger := otelslog.NewLogger("auth.middleware")

		sessionID := getSessionID(r)
		if sessionID == "" {
			authLogger.InfoContext(r.Context(), "Missing session", "path", r.URL.Path)
			http.Error(w, "Unauthorized: missing sessionID", http.StatusUnauthorized)
			return
		}

		client := getRedis()
		if client == nil {
			authLogger.ErrorContext(r.Context(), "Redis unavailable for auth")
			http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
			return
		}

		key := "session:" + sessionID
		_, err := client.Get(r.Context(), key).Result()
		if err != nil {
			authLogger.InfoContext(r.Context(), "Invalid or expired session", "path", r.URL.Path)
			http.Error(w, "Unauthorized: invalid or expired session", http.StatusUnauthorized)
			return
		}

		authLogger.InfoContext(r.Context(), "Auth OK", "path", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func TelemetryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracer.Start(r.Context(), r.URL.Path)
		defer span.End()
		r = r.WithContext(ctx) 

		startTime := time.Now()
		next.ServeHTTP(w, r)
		span.SetAttributes(attribute.String("method", r.Method))
		span.SetAttributes(attribute.String("path", r.URL.Path))
		span.SetAttributes(attribute.String("user_agent", r.UserAgent()))
		span.SetAttributes(attribute.String("referer", r.Referer()))
		span.SetAttributes(attribute.String("remote_addr", r.RemoteAddr))
		span.SetAttributes(attribute.String("duration", time.Since(startTime).String()))
	})
}
