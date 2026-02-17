package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

// Lista de usuarios en código: username -> password
var users = map[string]string{
	"admin": "secret",
	"demo":  "demo123",
}

func login(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	expectedPass, ok := users[creds.Username]
	if !ok || expectedPass != creds.Password {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	client := getRedis()
	if client == nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	sessionID := uuid.New().String()
	key := "session:" + sessionID
	ttl := parseSessionTTL()

	ctx := r.Context()
	if err := client.Set(ctx, key, creds.Username, ttl).Err(); err != nil {
		logger.ErrorContext(ctx, "Failed to store session", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"sessionID": sessionID})
}

func parseSessionTTL() time.Duration {
	s := os.Getenv("SESSION_TTL")
	if s == "" {
		return 24 * time.Hour
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return 24 * time.Hour
	}
	return d
}
