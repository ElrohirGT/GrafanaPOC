package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	// Handle SIGINT (CTRL+C) gracefully.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Set up OpenTelemetry.
	otelShutdown, err := setupOTelSDK(ctx)
	if err != nil {
		return err
	}
	// Handle shutdown properly so nothing leaks.
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	setupPyroscope()

	if err := initDB(ctx); err != nil {
		log.Printf("Warning: failed to init DB: %v", err)
	}
	defer func() {
		_ = closeDB(context.Background())
	}()

	if err := initRedis(ctx); err != nil {
		log.Printf("Warning: failed to init Redis: %v", err)
	}
	defer func() {
		_ = closeRedis(context.Background())
	}()

	// Start HTTP server.
	srv := &http.Server{
		Addr:         ":8080",
		BaseContext:  func(net.Listener) context.Context { return ctx },
		ReadTimeout:  time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      newHTTPHandler(),
	}
	srvErr := make(chan error, 1)
	go func() {
		srvErr <- srv.ListenAndServe()
	}()

	// Wait for interruption.
	select {
	case err = <-srvErr:
		// Error when starting HTTP server.
		return err
	case <-ctx.Done():
		// Wait for first CTRL+C.
		// Stop receiving signal notifications as soon as possible.
		stop()
	}

	// When Shutdown is called, ListenAndServe immediately returns ErrServerClosed.
	err = srv.Shutdown(context.Background())
	return err
}

func newHTTPHandler() http.Handler {
	mux := http.NewServeMux()

	handlerPublic := Chain(publicMiddlewares...)
	mux.Handle("/login", handlerPublic(http.HandlerFunc(login)))
	mux.Handle("/rolldice", handlerPublic(http.HandlerFunc(rolldice)))
	mux.Handle("/rolldice/{player}", handlerPublic(http.HandlerFunc(rolldice)))
	mux.Handle("/artists", handlerPublic(http.HandlerFunc(artists)))
	mux.Handle("/invoices/summary", handlerPublic(http.HandlerFunc(invoicesSummary)))

	handlerProtected := Chain(protectedMiddlewares...)
	mux.Handle("/heavy", handlerProtected(http.HandlerFunc(heavy)))
	mux.Handle("/heavy/{player}", handlerProtected(http.HandlerFunc(heavy)))
	mux.Handle("/stats", handlerProtected(http.HandlerFunc(stats)))

	// Add HTTP instrumentation for the whole server.
	handler := otelhttp.NewHandler(mux, "/")
	return handler
}
