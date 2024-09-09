# Build stage
FROM golang AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod
COPY go.mod ./

# Download all dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build for the host architecture
RUN CGO_ENABLED=0 go build -v -o go-weather-app

# Final stage
FROM debian:buster-slim

# Set the working directory inside the container
WORKDIR /root/

# Copy the binary from the build stage
COPY --from=builder /app/go-weather-app .

# Copy the templates directory
COPY --from=builder /app/templates ./templates

# Install ca-certificates for HTTPS requests
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./go-weather-app"]