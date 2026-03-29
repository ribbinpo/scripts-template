# ---------- BUILD STAGE ----------
  FROM golang:1.25-alpine AS builder

  WORKDIR /app
  
  # Install git (needed for go mod sometimes)
  RUN apk add --no-cache git
  
  # Copy go mod first (cache deps)
  COPY go.mod go.sum ./
  RUN go mod download
  
  # Copy source
  COPY . .
  
  # Build binary (static, small)
  RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
      go build -o service-b ./service-b
  
  # ---------- RUNTIME STAGE ----------
  FROM alpine:3.19
  
  WORKDIR /app
  
  # Install certs (for HTTPS calls)
  RUN apk add --no-cache ca-certificates
  
  # Copy binary only
  COPY --from=builder /app/service-b .
  
  # Expose port
  EXPOSE 4001
  
  # Run app
  CMD ["./service-b"]