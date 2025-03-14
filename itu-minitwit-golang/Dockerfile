FROM golang:1.24

LABEL org.opencontainers.image.source = "https://github.com/devops-q/devops"

# Set destination for COPY
WORKDIR /app

# Install necessary dependencies for CGO
RUN apt-get update && apt-get install -y --no-install-recommends gcc=9.3.0 && rm -rf /var/lib/apt/lists/*
# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build
RUN CGO_ENABLED=1 GOOS=linux go build -o ./itu-minitwit-golang ./cmd/server/main.go

EXPOSE 8080

# Run
CMD ["./itu-minitwit-golang"]