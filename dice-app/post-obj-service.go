package main

import (
	"context"
	"io"
	"net/http"
	"strings"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func PostObj(ctx context.Context) error {
	url := "https://api.restful-api.dev/objects"
	ctx, span := tracer.Start(ctx, "post-obj", trace.WithAttributes(
		attribute.String("url", url),
	))
	defer span.End()

	logger.InfoContext(ctx, "Adding obj to service...")
	reqBody := `
{
   "name": "Apple MacBook Pro 16",
   "data": {
      "year": 2019,
      "price": 1849.99,
      "CPU model": "Intel Core i9",
      "Hard disk size": "1 TB"
   }
}`
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(reqBody))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	logger.InfoContext(ctx, "Service responded!")
	if err != nil {
		span.AddEvent("Integration `post-obj`", trace.WithAttributes(
			attribute.Bool("success", false),
			attribute.String("request", reqBody),
			attribute.String("reason", err.Error()),
		))
		return err
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		span.AddEvent("Integration `post-obj`", trace.WithAttributes(
			attribute.Bool("success", false),
			attribute.String("request", reqBody),
			attribute.String("reason", err.Error()),
		))
		return err
	}

	span.AddEvent("Integration `post-obj`", trace.WithAttributes(
		attribute.Bool("success", true),
		attribute.String("request", reqBody),
		attribute.String("response", string(respBody)),
	))

	return nil
}
