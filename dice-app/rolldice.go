package main

import (
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"go.opentelemetry.io/contrib/bridges/otelslog"
)

const name = "go.opentelemetry.io/contrib/examples/dice"

var (
	tracer     = otel.Tracer(name)
	meter      = otel.Meter(name)
	logger     = otelslog.NewLogger(name)
	rollByUser metric.Int64Gauge
)

func init() {
	var err error

	rollByUser, err = meter.Int64Gauge("dice.rolls.byUser",
		metric.WithDescription("The rolls each user got"),
		metric.WithUnit("{roll}"),
	)
	if err != nil {
		panic(err)
	}
}

const rolldiceCacheTTL = 30 * time.Second

func rolldice(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "roll")
	defer span.End()

	playerName := r.PathValue("player")
	if playerName == "" {
		playerName = "Anonymous"
	}
	cacheKey := "rolldice:last:" + playerName

	roll := 1 + rand.Intn(6)
	rollValueAttr := attribute.Int("roll.value", roll)
	userAttr := attribute.String("user.name", playerName)

	// Cache roll in Redis
	if client := getRedis(); client != nil {
		_ = client.Set(ctx, cacheKey, roll, rolldiceCacheTTL).Err()
	}
	msg := playerName + " is rolling the dice"
	logger.InfoContext(ctx, msg, "result", roll)

	span.SetAttributes(rollValueAttr, userAttr)

	logger.InfoContext(ctx, "Getting first user...")
	err := GetUser(ctx)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to get user!", "error", err)
		span.RecordError(err)
		http.Error(w, "Failed to get user! "+err.Error(), http.StatusInternalServerError)
		return
	}

	logger.InfoContext(ctx, "Posting a custom obj...")
	err = PostObj(ctx)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to post obj!", "error", err)
		span.RecordError(err)
		http.Error(w, "Failed to post obj! "+err.Error(), http.StatusInternalServerError)
		return
	}

	rollByUser.Record(ctx, int64(roll), metric.WithAttributes(userAttr))
	resp := strconv.Itoa(roll) + "\n"
	if _, err := io.WriteString(w, resp); err != nil {
		logger.ErrorContext(ctx, "Write failed", "error", err)
	}
}
