# Build stage
FROM golang:1.22-alpine AS builder

ENV GOPROXY=https://goproxy.io,direct

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o simple-file-server .

# Runtime stage
FROM alpine:3.21

# Install ca-certificates for HTTPS requests if needed
RUN apk --no-cache add ca-certificates

RUN mkdir /data/

WORKDIR /

# Copy the binary from builder stage
COPY --from=builder /app/simple-file-server /

# Copy default config
COPY config.json .

# Create directories for data and temp
RUN mkdir -p /data /temp ./temp/MultiPart

# Expose the port
EXPOSE 60088

# Run the application
CMD ["/simple-file-server"]