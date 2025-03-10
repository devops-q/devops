# Build stage
FROM golang:1.24 AS builder

LABEL org.opencontainers.image.source = "https://github.com/devops-q/devops"

# Set destination for COPY
WORKDIR /app

# Install necessary dependencies for CGO
RUN apt-get update && apt-get install -y gcc

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o ./create_api_user ./cmd/create_api_user/main.go

# Final stage
FROM alpine:latest

# Set destination for COPY
WORKDIR /app

# Copy the binary from the build stage
COPY --from=builder /app/create_api_user .

# Set the entrypoint
ENTRYPOINT ./create_api_user -username=$API_USER -password=$API_PASSWORD