FROM golang:1.25.0-alpine AS builder

WORKDIR /app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Pass the service name as a build argument (default to screen-service)
ARG SERVICE="screen-service"

# Build the Go binary natively
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/service ./cmd/${SERVICE}/main.go

# Create the minimal runtime image
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app
COPY --from=builder /bin/service /app/service

# Expose the standard port (can be overridden by docker-compose)
EXPOSE 50051 50052

ENTRYPOINT ["/app/service"]
