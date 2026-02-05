package main

import (
	"context"
	"errors"
	"io"
	"math/rand"
	"net/http"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var ErrSimulatedError = errors.New("emulating: something went wrong during parsing of body")

func GetUser(ctx context.Context) error {
	url := "https://randomuser.me/api/"
	ctx, span := tracer.Start(ctx, "get-user", trace.WithAttributes(
		attribute.String("url", url),
	))
	defer span.End()

	logger.InfoContext(ctx, "Getting user from api...")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	logger.InfoContext(ctx, "Service responded!")
	if err != nil {
		span.AddEvent("Integration `get-user` failed!", trace.WithAttributes(
			attribute.String("reason", err.Error()),
		))
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		span.AddEvent("Integration `get-user` failed!", trace.WithAttributes(
			attribute.String("reason", err.Error()),
		))
		return err
	}
	if rand.Float32() < 0.1 {
		span.AddEvent("Integration `get-user` failed!", trace.WithAttributes(
			attribute.String("reason", ErrSimulatedError.Error()),
		))
		return ErrSimulatedError
	}

	span.AddEvent("Integration `get-user` success!", trace.WithAttributes(
		attribute.String("body", string(body)),
	))

	return nil
}
