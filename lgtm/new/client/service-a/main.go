package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type PingPongResp struct {
	Message string `json:"message"`
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
		ctx := context.WithValue(c.Request.Context(), "request_id", requestId)
		c.Request = c.Request.WithContext(ctx)

		// ----- process request -----
		c.Next()

		// ----- duration -----
		duration := time.Since(t)

		fields := []zap.Field{
			zap.String("service", "service-a"),
			zap.String("request_id", requestId),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Int64("duration_ms", duration.Milliseconds()),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.String("client_ip", c.ClientIP()),
		}

		fields = append(fields,
			zap.String("trace_id", ""),
			zap.String("span_id", ""),
		)

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

	router.GET("/ping-service-b", func(c *gin.Context) {
		req, err := http.NewRequest("GET", "http://service-b:4001/ping", nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "error creating request",
			})
			return
		}
		req.Header.Set("x-request-id", c.GetHeader("x-request-id"))

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
			"message": "message from service-b: " + responseBody.Message,
		})
	})

	router.Run(":4000")
}
