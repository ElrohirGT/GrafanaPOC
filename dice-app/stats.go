package main

import (
	"encoding/json"
	"net/http"
	"sync"
)

type statsResponse struct {
	ArtistsCount  int     `json:"artistsCount"`
	TracksCount   int     `json:"tracksCount"`
	RedisKeys     int64   `json:"redisKeys,omitempty"`
	GetUserOK     bool    `json:"getUserOk"`
	PostObjOK     bool    `json:"postObjOk"`
}

func stats(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "stats")
	defer span.End()

	resp := statsResponse{}

	var wg sync.WaitGroup

	// SQL: count artists
	wg.Add(1)
	go func() {
		defer wg.Done()
		db := getDB()
		if db != nil {
			var n int
			if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM artists").Scan(&n); err == nil {
				resp.ArtistsCount = n
			}
		}
	}()

	// SQL: count tracks
	wg.Add(1)
	go func() {
		defer wg.Done()
		db := getDB()
		if db != nil {
			var n int
			if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM tracks").Scan(&n); err == nil {
				resp.TracksCount = n
			}
		}
	}()

	// Redis: DBSize
	wg.Add(1)
	go func() {
		defer wg.Done()
		client := getRedis()
		if client != nil {
			if n, err := client.DBSize(ctx).Result(); err == nil {
				resp.RedisKeys = n
			}
		}
	}()

	// GetUser
	wg.Add(1)
	go func() {
		defer wg.Done()
		resp.GetUserOK = (GetUser(ctx) == nil)
	}()

	// PostObj
	wg.Add(1)
	go func() {
		defer wg.Done()
		resp.PostObjOK = (PostObj(ctx) == nil)
	}()

	wg.Wait()

	data, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(data)
}
