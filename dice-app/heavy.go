package main

import (
	"context"
	"io"
	"net/http"
	"sync"

	"go.opentelemetry.io/otel/attribute"
)

// const name = "go.opentelemetry.io/contrib/examples/dice"
//
// var (
//
//	tracer     = otel.Tracer(name)
//	meter      = otel.Meter(name)
//	logger     = otelslog.NewLogger(name)
//	rollByUser metric.Int64Gauge
//
// )
//
//	func init() {
//		var err error
//
//		rollByUser, err = meter.Int64Gauge("dice.rolls.byUser",
//			metric.WithDescription("The rolls each user got"),
//			metric.WithUnit("{roll}"),
//		)
//		if err != nil {
//			panic(err)
//		}
//	}

// Heavy handler simulation
func heavy(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "heavy")
	defer span.End()

	playerName := r.PathValue("player")
	if playerName == "" {
		playerName = "Anonymous"
	}
	userAttr := attribute.String("user.name", playerName)
	span.SetAttributes(userAttr)

	logger.InfoContext(ctx, "Simulating CPU intensive work...")
	func() {
		_, span := tracer.Start(ctx, "cpuLoop")
		defer span.End()
		for i := range 1_000_000_000 {
			if i%100_000_000 == 0 {
				logger.InfoContext(ctx, "iteration checkpoint", "iteration", i)
			}
		}
	}()

	logger.InfoContext(ctx, "Simulating RAM intensive work...")
	func() {
		_, span := tracer.Start(ctx, "ramConsumption")
		defer span.End()
		a := make([]uint8, 1_000_000) // Allocate 1Mb of data

		// Set some random values
		a[5000] = 1
		a[999999] = 4
	}()

	workersCtx, cancelWorkersCtx := context.WithCancelCause(ctx)
	func() {
		defer cancelWorkersCtx(nil)
		ctx, span := tracer.Start(workersCtx, "goRoutines")
		defer span.End()

		logger.InfoContext(ctx, "Simulate Go routines...")
		var wg sync.WaitGroup
		defer wg.Wait()

		wg.Go(func() {
			logger.InfoContext(ctx, "Getting first user...")
			err := GetUser(ctx)
			if err != nil {
				logger.ErrorContext(ctx, "Failed to get user!", "error", err)
				span.RecordError(err)
				cancelWorkersCtx(err)
				return
			}
		})

		wg.Go(func() {
			logger.InfoContext(ctx, "Posting a custom obj...")
			err := PostObj(ctx)
			if err != nil {
				logger.ErrorContext(ctx, "Failed to post obj!", "error", err)
				span.RecordError(err)
				cancelWorkersCtx(err)
				return
			}
		})
	}()

	if err := context.Cause(workersCtx); err != nil && err != context.Canceled {
		logger.InfoContext(ctx, "Workers failed!")
		http.Error(w, "Workers failed! "+err.Error(), 500)
		return
	}

	logger.InfoContext(ctx, "Done!")
	_, err := io.WriteString(w, "OK!")
	if err != nil {
		logger.ErrorContext(ctx, "Failed to write response", err)
	}
}
