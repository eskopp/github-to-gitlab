# Use the official Go image to build the app
FROM golang:1.23.2 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o /app/main .

# Use Artix Linux as the base for running the app
FROM artixlinux/base:latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Install Git on Artix
RUN pacman -Sy --noconfirm git

# Copy the pre-built binary from the builder stage
COPY --from=builder /app/main .

# Run the Go binary
CMD ["/main"]
