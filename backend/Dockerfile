FROM golang:1.23-alpine

# Set working directory
WORKDIR /app

# Copy go.mod + go.sum then download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Grab the rest of the application code
COPY . .

# Build the Go binary
RUN go build -o api ./cmd/api

# Exposes 8000 port for API
EXPOSE 8000

# Run the API
CMD ["./api"]