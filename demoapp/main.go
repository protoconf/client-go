package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/avast/retry-go"
	protoconfLoader "github.com/protoconf/client-go"
	democonfig "github.com/protoconf/client-go/demoapp/config/src/demo/v1"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var ConfigDir = flag.String("config-dir", "local", "directory to load config from")

// main function that parses flags, sets up logger, telemetry, loads config, and starts an HTTP server.
func main() {
	flag.Parse()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	logger := setupLogger()
	tracerProvider, meterProvider := setupTelemetry(ctx)
	defer shutdownTelemetry(ctx, tracerProvider, meterProvider)

	slog.Info("starting")
	config := loadConfig(ctx, logger)
	startHTTPServer(ctx, logger, config)
}

func setupLogger() *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, nil)
	return slog.New(handler)
}

// setupTelemetry initializes and configures OpenTelemetry tracing and metrics providers.
// It creates gRPC-based exporters for tracing and metrics, sets up the TracerProvider and MeterProvider,
// and returns them for further use.
//
// Parameters:
// ctx (context.Context): The context for the setup process.
//
// Returns:
// tracerProvider (*trace.TracerProvider): The configured TracerProvider for tracing.
// meterProvider (*metric.MeterProvider): The configured MeterProvider for metrics.
func setupTelemetry(ctx context.Context) (*trace.TracerProvider, *metric.MeterProvider) {
	// Create a gRPC-based exporter for tracing
	expTracer, err := otlptracegrpc.New(ctx)
	if err != nil {
		slog.Error("Error in setupTelemetry", "error", err)
		return nil, nil
	}

	// Set up the TracerProvider with the created exporter
	tracerProvider := trace.NewTracerProvider(trace.WithBatcher(expTracer))
	otel.SetTracerProvider(tracerProvider)

	// Create a gRPC-based exporter for metrics
	expMeter, err := otlpmetricgrpc.New(ctx)
	if err != nil {
		panic(err)
	}

	// Set up the MeterProvider with the created exporter
	meterProvider := metric.NewMeterProvider(metric.WithReader(metric.NewPeriodicReader(expMeter)))
	otel.SetMeterProvider(meterProvider)

	return tracerProvider, meterProvider
}

// shutdownTelemetry is responsible for gracefully shutting down the OpenTelemetry tracing and metrics providers.
// It takes a context, the TracerProvider, and the MeterProvider as parameters.
// It shuts down the TracerProvider and the MeterProvider by calling their respective Shutdown methods.
// If an error occurs during the shutdown process, it panics.
//
// Parameters:
// ctx (context.Context): The context for the shutdown process.
// tracerProvider (*trace.TracerProvider): The TracerProvider to be shut down.
// meterProvider (*metric.MeterProvider): The MeterProvider to be shut down.
func shutdownTelemetry(ctx context.Context, tracerProvider *trace.TracerProvider, meterProvider *metric.MeterProvider) {
	// Shutdown the TracerProvider
	if err := tracerProvider.Shutdown(ctx); err != nil {
		panic(err)
	}

	// Shutdown the MeterProvider
	if err := meterProvider.Shutdown(ctx); err != nil {
		panic(err)
	}
}

// loadConfig loads the configuration from the specified directory and returns a DemoConfig object.
// It also sets up a watcher to monitor for changes in the configuration file and logs any errors encountered.
//
// Parameters:
// ctx (context.Context): The context for the configuration loading process.
// logger (*slog.Logger): The logger to be used for logging.
//
// Returns:
// *democonfig.DemoConfig: The loaded DemoConfig object.
func loadConfig(ctx context.Context, logger *slog.Logger) *democonfig.DemoConfig {
	// Initialize the DemoConfig object with a default title.
	config := &democonfig.DemoConfig{Title: "dummy"}
	configHolder := &atomic.Pointer[democonfig.DemoConfig]{}

	// Get the hostname of the pod.
	podName, err := os.Hostname()
	if err == nil {
		// Set the pod name in the DemoConfig object if the hostname is obtained successfully.
		config.PodName = podName
	}
	configHolder.Store(config)
	// Create a new configuration loader using protoconfLoader.
	configLib, err := protoconfLoader.NewConfiguration(config, "demoapp", protoconfLoader.WithLogger(logger))
	if err != nil {
		// Log an error if the configuration loader cannot be created.
		logger.Error("failed to create protoconf loader", "error", err)
	}

	// Load the configuration from the specified directory.
	err = configLib.LoadConfig(*ConfigDir, "config.json")
	if err != nil {
		// Log an error if the configuration cannot be loaded.
		logger.Error("failed to load config", "error", err)
	}

	// Set up a callback function to be executed when the configuration changes.
	configLib.OnConfigChange(func(p proto.Message) {
		// Log the new configuration.
		logger.Info("got new config", "config", p)
	})

	// Retry the WatchConfig function in case of errors.
	err = retry.Do(func() error {
		// Watch for changes in the configuration file.
		err := configLib.WatchConfig(ctx)
		if err != nil {
			// Log an error if the WatchConfig function fails.
			logger.Error("retry attempt failed", "error", err)
		}
		return err
	})
	if err != nil {
		// Log an error if the WatchConfig function fails after all retry attempts.
		logger.Error("failed to watch config", "error", err)
	}

	// Return the loaded DemoConfig object.
	return configHolder.Load()
}

// startHTTPServer sets up and starts an HTTP server that listens on port 8181.
// It serves static files from the specified configuration directory and provides an endpoint
// for retrieving the current configuration. The server is instrumented with OpenTelemetry for tracing.
//
// Parameters:
// ctx (context.Context): The context for the server.
// logger (*slog.Logger): The logger to be used for logging.
// config (*democonfig.DemoConfig): The current configuration.
func startHTTPServer(ctx context.Context, logger *slog.Logger, config *democonfig.DemoConfig) {
	// Initialize a variable to keep track of the total number of requests.
	var totalRequests uint64

	// Create a new HTTP server with specified timeouts.
	srv := &http.Server{
		Addr:         ":8181",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// Register a function to be executed when the server is shutting down.
	srv.RegisterOnShutdown(func() {
		_ = srv.Shutdown(ctx)
	})

	// Create a new ServeMux for routing requests.
	mux := http.NewServeMux()

	// Handle requests for static files from the configuration directory.
	mux.Handle("/configs/", http.StripPrefix("/configs/", http.FileServer(http.Dir(*ConfigDir))))

	// Handle requests for the root path.
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Increment the total number of requests.
		atomic.AddUint64(&totalRequests, 1)

		// Log a debug message for each request.
		logger.Debug("got request")

		// Add a custom header to the response with the title from the configuration.
		w.Header().Add("x-demo-title", config.GetTitle())

		// Update the total number of requests in the configuration.
		config.TotalRequests = totalRequests

		// Marshal the configuration to JSON.
		b, err := protojson.MarshalOptions{}.Marshal(config)
		if err != nil {
			// If there's an error, return a 500 Internal Server Error response.
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "error: %v", err)
			return
		}

		// Return a 200 OK response with the JSON-marshaled configuration.
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	})

	// Instrument the ServeMux with OpenTelemetry for tracing.
	srv.Handler = otelhttp.NewHandler(mux, "get")

	// Start the server and log any errors that occur.
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("server error", "error", err)
	}
}
