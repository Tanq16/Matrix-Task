# Build stage
FROM golang:alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/build/matrix /app/cmd/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy only the binary
COPY --from=builder /app/build/matrix .

EXPOSE 8080

CMD ["./matrix"]
