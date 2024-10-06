# Use the official Go image (version 1.23.2) to build the app
FROM golang:1.23.2-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o main .

# Use a minimal Alpine image for running the application
FROM alpine:3.20

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the pre-built binary from the builder stage
COPY --from=builder /app/main .

# Install Git to enable Git commands in the container
RUN apk add --no-cache git

# Expose any ports (if needed)
# EXPOSE 8080

# Command to run the executable
ENTRYPOINT ["./main"]
