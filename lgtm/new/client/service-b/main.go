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
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
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

func InitTracer(prop propagation.TextMapPropagator) func(context.Context) error {
	ctx := context.Background()

	exp, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint("lgtm-tempo:4318"),
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

	prop := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)

	shutdown := InitTracer(prop)
	defer shutdown(context.Background())

	router.Use(otelgin.Middleware("service-b", otelgin.WithPropagators(prop)))

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
