FROM golang:1.20-alpine

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum then download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go binary
RUN go build -o api ./cmd/api

# Expose the port your API listens on (adjust if necessary)
EXPOSE 8000

# Run the API
CMD ["./api"]