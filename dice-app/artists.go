package main

import (
	"encoding/json"
	"net/http"
	"time"
)

const artistsCacheTTL = 45 * time.Second
const artistsCacheKey = "artists:list"

type artist struct {
	ArtistID int    `json:"artistId"`
	Name     string `json:"name"`
}

func artists(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "artists")
	defer span.End()

	client := getRedis()
	if client != nil {
		cached, err := client.Get(ctx, artistsCacheKey).Bytes()
		if err == nil {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(cached)
			return
		}
	}

	database := getDB()
	if database == nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	rows, err := database.QueryContext(ctx, "SELECT ArtistId, Name FROM artists LIMIT 10")
	if err != nil {
		logger.ErrorContext(ctx, "Query artists failed", "error", err)
		span.RecordError(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var list []artist
	for rows.Next() {
		var a artist
		if err := rows.Scan(&a.ArtistID, &a.Name); err != nil {
			logger.ErrorContext(ctx, "Scan artist failed", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		list = append(list, a)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(list)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if client != nil {
		_ = client.Set(ctx, artistsCacheKey, data, artistsCacheTTL).Err()
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(data)
}

