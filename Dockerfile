FROM golang:alpine

# Install git, bash, and build essentials (if needed by pty/CGO)
RUN apk add --no-cache git bash build-base

WORKDIR /app

# Copy mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application
COPY . .

# Build the main server application
RUN go build -o main ./cmd/server/main.go

# Expose port 8080
EXPOSE 8080

# Run the binary
CMD ["./main"]
