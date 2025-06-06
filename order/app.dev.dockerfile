FROM golang:1.24.3-alpine3.22

WORKDIR /app

# Install Air for hot reload
RUN go install github.com/air-verse/air@latest

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code (will be volume mounted in dev)
COPY . .

# Use Air for hot reload
CMD ["air", "-c", ".air.order.toml"] 