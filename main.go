# Use the official Go image to build the app
FROM golang:1.19-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod ./
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o main .

# Use a minimal image for running the application
FROM alpine:3.15

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the pre-built binary from the builder stage
COPY --from=builder /app/main .

# Copy Git to be able to use git commands
RUN apk add --no-cache git

# Expose any ports (if needed)
# EXPOSE 8080

# Command to run the executable
ENTRYPOINT ["./main"]
