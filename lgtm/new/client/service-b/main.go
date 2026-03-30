package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	metricapi "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ctxKey int

const requestIDKey ctxKey = iota

type PingPongResp struct {
	Message string `json:"message"`
}

type HTTPMetrics struct {
	requestCounter metricapi.Int64Counter
	latencyMs      metricapi.Float64Histogram
}

func InitMetric() (func(context.Context) error, *HTTPMetrics) {
	ctx := context.Background()

	exp, err := otlpmetrichttp.New(ctx,
		otlpmetrichttp.WithEndpoint("alloy:4318"),
		otlpmetrichttp.WithInsecure(),
	)
	if err != nil {
		log.Fatal(err)
	}

	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exp)),
		sdkmetric.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("service-b"),
			semconv.DeploymentEnvironmentName("dev"),
		)),
	)

	otel.SetMeterProvider(provider)
	meter := provider.Meter("service-b")

	requestCounter, err := meter.Int64Counter(
		"http.server.request.count",
		metricapi.WithDescription("Total HTTP requests handled by service-b"),
	)
	if err != nil {
		log.Fatal(err)
	}
	latencyMs, err := meter.Float64Histogram(
		"http.server.request.duration.ms",
		metricapi.WithDescription("HTTP request duration in milliseconds for service-b"),
	)
	if err != nil {
		log.Fatal(err)
	}

	return provider.Shutdown, &HTTPMetrics{
		requestCounter: requestCounter,
		latencyMs:      latencyMs,
	}
}

func MetricsMiddleware(metrics *HTTPMetrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		attrs := []attribute.KeyValue{
			attribute.String("service.name", "service-b"),
			attribute.String("http.method", c.Request.Method),
			attribute.String("http.route", c.FullPath()),
			attribute.Int("http.status_code", c.Writer.Status()),
		}

		ctx := c.Request.Context()
		metrics.requestCounter.Add(ctx, 1, metricapi.WithAttributes(attrs...))
		metrics.latencyMs.Record(
			ctx,
			float64(time.Since(start).Milliseconds()),
			metricapi.WithAttributes(attrs...),
		)
	}
}

func InitTracer(prop propagation.TextMapPropagator) func(context.Context) error {
	ctx := context.Background()

	exp, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint("tempo:4318"),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		log.Fatal(err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("service-b"),
			semconv.DeploymentEnvironmentName("dev"),
		)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(prop)

	return tp.Shutdown
}

func LoggerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		// ----- request_id -----
		requestId := c.GetHeader("x-request-id")
		if requestId == "" {
			requestId = uuid.New().String()
		}
		// set for downstream + response
		c.Request.Header.Set("x-request-id", requestId)
		c.Writer.Header().Set("x-request-id", requestId)

		// store in context
		ctx := context.WithValue(c.Request.Context(), requestIDKey, requestId)
		c.Request = c.Request.WithContext(ctx)

		// ----- process request -----
		c.Next()

		// ----- duration -----
		duration := time.Since(t)

		fields := []zap.Field{
			zap.String("service", "service-b"),
			zap.String("request_id", requestId),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Int64("duration_ms", duration.Milliseconds()),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.String("client_ip", c.ClientIP()),
		}

		span := trace.SpanFromContext(c.Request.Context())
		sc := span.SpanContext()

		if sc.IsValid() {
			fields = append(fields,
				zap.String("trace_id", sc.TraceID().String()),
				zap.String("span_id", sc.SpanID().String()),
			)
		}

		// ----- error handling -----
		if len(c.Errors) > 0 {
			fields = append(fields, zap.Error(c.Errors[0]))
			logger.Error("request failed", fields...)
		} else {
			logger.Info("request completed", fields...)
		}
	}
}

func main() {
	router := gin.Default()
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	shutdownMetric, httpMetrics := InitMetric()
	defer shutdownMetric(context.Background())

	prop := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)

	shutdown := InitTracer(prop)
	defer shutdown(context.Background())

	router.Use(otelgin.Middleware("service-b", otelgin.WithPropagators(prop)))
	router.Use(MetricsMiddleware(httpMetrics))

	cfg := zap.NewProductionConfig()

	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.LevelKey = "level"
	cfg.EncoderConfig.MessageKey = "message"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, _ := cfg.Build()
	defer logger.Sync()

	router.Use(LoggerMiddleware(logger))

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.GET("/ping-service-a", func(c *gin.Context) {
		req, err := http.NewRequestWithContext(c.Request.Context(), "GET", "http://service-a:4000/ping", nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "error creating request",
			})
			return
		}
		req.Header.Set("x-request-id", c.GetHeader("x-request-id"))
		otel.GetTextMapPropagator().Inject(c.Request.Context(), propagation.HeaderCarrier(req.Header))

		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "error doing request",
			})
			return
		}
		defer resp.Body.Close()

		var responseBody PingPongResp
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "error reading response body",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "message from service-a: " + responseBody.Message,
		})
	})

	router.Run(":4001")
}
