# Build stage
FROM golang:1.24 AS builder

LABEL org.opencontainers.image.source="https://github.com/devops-q/devops"

# Set destination for COPY
WORKDIR /app

# Install necessary dependencies for CGO
RUN apt-get update && apt-get install -y --no-install-recommends gcc=4:12.2.0-3

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o ./itu-minitwit-golang ./cmd/server/main.go

# Final stage
FROM scratch

# Set destination for COPY
WORKDIR /app

COPY --from=builder /app/web ./web

# Copy the binary from the build stage
COPY --from=builder /app/itu-minitwit-golang .

EXPOSE 8080

# Run
CMD ["./itu-minitwit-golang"]