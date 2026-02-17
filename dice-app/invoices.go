package main

import (
	"encoding/json"
	"net/http"
	"time"
)

const invoicesCacheTTL = 45 * time.Second
const invoicesCacheKey = "invoices:summary"

type invoiceSummary struct {
	BillingCountry string  `json:"billingCountry"`
	TotalSales     float64 `json:"totalSales"`
}

func invoicesSummary(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "invoices-summary")
	defer span.End()

	client := getRedis()
	if client != nil {
		cached, err := client.Get(ctx, invoicesCacheKey).Bytes()
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

	rows, err := database.QueryContext(ctx,
		"SELECT BillingCountry, SUM(Total) AS TotalSales FROM invoices GROUP BY BillingCountry")
	if err != nil {
		logger.ErrorContext(ctx, "Query invoices failed", "error", err)
		span.RecordError(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var list []invoiceSummary
	for rows.Next() {
		var inv invoiceSummary
		if err := rows.Scan(&inv.BillingCountry, &inv.TotalSales); err != nil {
			logger.ErrorContext(ctx, "Scan invoice failed", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		list = append(list, inv)
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
		_ = client.Set(ctx, invoicesCacheKey, data, invoicesCacheTTL).Err()
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(data)
}
