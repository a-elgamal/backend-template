package telemetry

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-logr/logr"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// Middleware returns middleware that will trace incoming requests.
// The service parameter should describe the name of the (virtual)
// server handling the request.
func Middleware(logger logr.Logger) gin.HandlerFunc {
	// Create the metrics
	meter := otel.Meter("internal/telemetry")
	attemptsCounter, _ := meter.Int64Counter("http.request_count", metric.WithDescription("Number of Requests"), metric.WithUnit("Count"))
	totalDuration, _ := meter.Int64Histogram("http.duration", metric.WithDescription("Time Taken by request"), metric.WithUnit("Milliseconds"))
	activeRequestsCounter, _ := meter.Int64UpDownCounter("http.active_requests", metric.WithDescription("Number of requests inflight"), metric.WithUnit("Count"))
	requestSize, _ := meter.Int64Histogram("http.request_content_length", metric.WithDescription("Request Size"), metric.WithUnit("Bytes"))
	responseSize, _ := meter.Int64Histogram("http.response_content_length", metric.WithDescription("Response Size"), metric.WithUnit("Bytes"))

	return func(ginCtx *gin.Context) {
		start := time.Now()

		// Adding logger into the request context
		ctx := logr.NewContext(ginCtx.Request.Context(), logger)
		ginCtx.Request = ginCtx.Request.WithContext(ctx)

		attrs := []attribute.KeyValue{
			semconv.HTTPRequestMethodKey.String(ginCtx.Request.Method),
		}

		route := ginCtx.FullPath()
		if route != "" {
			attrs = append(attrs, semconv.HTTPRouteKey.String(route))
		}

		activeRequestsCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
		defer activeRequestsCounter.Add(ctx, -1, metric.WithAttributes(attrs...))

		defer func() {

			resAttributes := append(attrs, semconv.HTTPResponseStatusCodeKey.Int(ginCtx.Writer.Status()))

			attemptsCounter.Add(ctx, 1, metric.WithAttributes(resAttributes...))

			rqSize := computeApproximateRequestSize(ginCtx.Request)
			requestSize.Record(ctx, rqSize, metric.WithAttributes(resAttributes...))
			responseSize.Record(ctx, int64(ginCtx.Writer.Size()), metric.WithAttributes(resAttributes...))

			totalDuration.Record(ctx, time.Since(start).Milliseconds(), metric.WithAttributes(resAttributes...))
		}()

		ginCtx.Next()

		if ginCtx.Writer.Status() >= 300 {
			logger.V(1).Info("Non-2xx status response detected!", "code", ginCtx.Writer.Status(), "url", ginCtx.Request.URL)
		}
	}
}

func computeApproximateRequestSize(r *http.Request) int64 {
	s := 0
	if r.URL != nil {
		s = len(r.URL.Path)
	}

	s += len(r.Method)
	s += len(r.Proto)
	for name, values := range r.Header {
		s += len(name)
		for _, value := range values {
			s += len(value)
		}
	}
	s += len(r.Host)

	// N.B. r.Form and r.MultipartForm are assumed to be included in r.URL.

	if r.ContentLength != -1 {
		s += int(r.ContentLength)
	}
	return int64(s)
}
