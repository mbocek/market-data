# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /market-data ./cmd/market-data

# Final stage
FROM alpine:3.18

WORKDIR /

# Copy the binary from builder
COPY --from=builder /market-data /market-data

# Expose port
EXPOSE 8080

# Run the binary
CMD ["/market-data"]
