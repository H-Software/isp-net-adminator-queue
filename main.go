package main

import (
	"context"
	"github.com/hibiken/asynq"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"os"
	// "net/http"
	"github.com/h-software/isp-net-adminator-queue/internal/flag"
	"github.com/h-software/isp-net-adminator-queue/internal/log"
	"github.com/h-software/isp-net-adminator-queue/internal/queue"
	"sync"
	"time"
	// "github.com/gorilla/handlers"
	// "github.com/gorilla/mux"
	// "go.opentelemetry.io/otel/metric"
	// "github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// Logger supports logging at various log levels.
type Logger interface {
	// Debug logs a message at Debug level.
	Debug(args ...interface{})

	// Info logs a message at Info level.
	Info(args ...interface{})

	// Warn logs a message at Warning level.
	Warn(args ...interface{})

	// Error logs a message at Error level.
	Error(args ...interface{})

	// Fatal logs a message at Fatal level
	// and process will exit with status set to 1.
	Fatal(args ...interface{})
}

type Config struct {
	// ErrorHandler handles errors returned by the task handler.
	//
	// HandleError is invoked only if the task handler returns a non-nil error.
	//
	// Example:
	//
	//     func reportError(ctx context, task *asynq.Task, err error) {
	//         retried, _ := asynq.GetRetryCount(ctx)
	//         maxRetry, _ := asynq.GetMaxRetry(ctx)
	//     	   if retried >= maxRetry {
	//             err = fmt.Errorf("retry exhausted for task %s: %w", task.Type, err)
	//     	   }
	//         errorReportingService.Notify(err)
	//     })
	//
	//     ErrorHandler: asynq.ErrorHandlerFunc(reportError)
	//
	//    we can also handle panic error like:
	//     func reportError(ctx context, task *asynq.Task, err error) {
	//         if asynq.IsPanicError(err) {
	//	          errorReportingService.Notify(err)
	// 	       }
	//     })
	//
	//     ErrorHandler: asynq.ErrorHandlerFunc(reportError)
	ErrorHandler ErrorHandler

	// Logger specifies the logger used by the server instance.
	//
	// If unset, default logger is used.
	Logger Logger

	// LogLevel specifies the minimum log level to enable.
	//
	// If unset, InfoLevel is used by default.
	LogLevel LogLevel
}

// LogLevel represents logging level.
//
// It satisfies flag.Value interface.
type LogLevel int32

// An ErrorHandler handles an error occurred during task processing.
type ErrorHandler interface {
	HandleError(ctx context.Context, err error)
}

// The ErrorHandlerFunc type is an adapter to allow the use of  ordinary functions as a ErrorHandler.
// If f is a function with the appropriate signature, ErrorHandlerFunc(f) is a ErrorHandler that calls f.
type ErrorHandlerFunc func(ctx context.Context, err error)

// HandleError calls fn(ctx, task, err)
func (fn ErrorHandlerFunc) HandleError(ctx context.Context, err error) {
	fn(ctx, err)
}

var (
	logger            *log.Logger
	resource          *sdkresource.Resource
	initResourcesOnce sync.Once
)

const (
	DEFAULT_RELOAD_INTERVAL = 10

	// httpAddress        = ":8080"
	// httpMetricsAddress = ":8081"

	// meterName = "adminator-queue-prometheus"
)

func init() {
	logger = log.NewLogger(nil)
}

func initResource() *sdkresource.Resource {
	initResourcesOnce.Do(func() {
		extraResources, _ := sdkresource.New(
			context.Background(),
			sdkresource.WithOS(),
			sdkresource.WithProcess(),
			sdkresource.WithContainer(),
			sdkresource.WithHost(),
		)
		resource, _ = sdkresource.Merge(
			sdkresource.Default(),
			extraResources,
		)
	})
	return resource
}

func initTracerProvider() *sdktrace.TracerProvider {
	ctx := context.Background()

	exporter, err := otlptracegrpc.New(ctx)
	if err != nil {
		logger.Fatalf("OTLP Trace gRPC Creation: %v", err)
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(initResource()),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp
}

func initMeterProvider() *sdkmetric.MeterProvider {
	ctx := context.Background()

	exporter, err := otlpmetricgrpc.New(ctx)
	if err != nil {
		logger.Fatalf("new otlp metric grpc exporter failed: %v", err)
	}

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter)),
		sdkmetric.WithResource(initResource()),
	)
	otel.SetMeterProvider(mp)
	return mp
}

// func serveMetrics() {
// 	log.Printf("serving metrics at %s/metrics", httpMetricsAddress)
// 	http.Handle("/metrics", promhttp.Handler())

// 	s := &http.Server{
// 		Addr:           httpMetricsAddress,
// 		Handler:        nil,
// 		ReadTimeout:    10 * time.Second,
// 		WriteTimeout:   10 * time.Second,
// 	}

// 	err := s.ListenAndServe()
// 	if err != nil {
// 		log.Printf("metrics: error serving http: %v", err)
// 		return
// 	}
// }

func main() {

	if flag.FlagHelp {
		flag.PrintDefaults()
		os.Exit(0)
	}

	tp := initTracerProvider()
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			logger.Infof("Tracer Provider Shutdown: %v", err)
		}
		logger.Infof("Shutdown tracer provider")
	}()

	mp := initMeterProvider()
	defer func() {
		if err := mp.Shutdown(context.Background()); err != nil {
			logger.Fatalf("Error shutting down meter provider: %v", err)
		}
		logger.Infof("Shutdown meter provider")
	}()

	err := runtime.Start(runtime.WithMinimumReadMemStatsInterval(time.Second))
	if err != nil {
		logger.Fatal(err)
	}

	// Initialize router
	// router := mux.NewRouter().StrictSlash(true)

	// Add apache-like logging to all routes
	// loggedRouter := handlers.CombinedLoggingHandler(os.Stdout, router)

	// log.Infof("adminator-queue server started on port: %s", httpAddress)

	// s := &http.Server{
	// 	Addr:         httpAddress,
	// 	Handler:      nil,
	// 	ReadTimeout:  10 * time.Second,
	// 	WriteTimeout: 10 * time.Second,
	// }

	// // start server
	// err = s.ListenAndServe()
	// if err != nil {
	// 	panic(err)
	// }

	srv := queue.RunServer()

	amux := asynq.NewServeMux()
	amux.HandleFunc(queue.TypeAdminatorWorkItemAgg, queue.HandleWorkItemAggTask)

	if err := srv.Run(amux); err != nil {
		logger.Fatalf("Failed to start the server: %v", err)
	}
}
