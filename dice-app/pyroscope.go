package main

import (
	"os"
	"runtime"

	otelpyroscope "github.com/grafana/otel-profiling-go"
	"github.com/grafana/pyroscope-go"
	"go.opentelemetry.io/otel"
)

func setupPyroscope() {
	// These 2 lines are only required if you're using mutex or block profiling
	// Read the explanation below for how to set these rates:
	runtime.SetMutexProfileFraction(5)
	runtime.SetBlockProfileRate(5)

	// Wrap it with otelpyroscope tracer provider.
	otel.SetTracerProvider(otelpyroscope.NewTracerProvider(otel.GetTracerProvider()))

	pyroscope.Start(pyroscope.Config{
		// ApplicationName: "dice.golang.app",
		ApplicationName: os.Getenv("APP_NAME"),

		// replace this with the address of pyroscope server
		// ServerAddress: "http://pyroscope-server:4040",
		ServerAddress: os.Getenv("PYROSCOPE_HOST"),

		// you can disable logging by setting this to nil
		Logger: pyroscope.StandardLogger,

		// you can provide static tags via a map:
		Tags: map[string]string{"hostname": os.Getenv("HOSTNAME")},

		ProfileTypes: []pyroscope.ProfileType{
			// these profile types are enabled by default:
			pyroscope.ProfileCPU,
			pyroscope.ProfileAllocObjects,
			pyroscope.ProfileAllocSpace,
			pyroscope.ProfileInuseObjects,
			pyroscope.ProfileInuseSpace,

			// these profile types are optional:
			pyroscope.ProfileGoroutines,
			pyroscope.ProfileMutexCount,
			pyroscope.ProfileMutexDuration,
			pyroscope.ProfileBlockCount,
			pyroscope.ProfileBlockDuration,
		},
	})

}
