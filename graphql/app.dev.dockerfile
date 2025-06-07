FROM golang:1.24.3-alpine3.22 AS dev

RUN apk update && apk --no-cache add build-base ca-certificates
WORKDIR /app

# Install Air for hot reloading
RUN go install github.com/air-verse/air@latest

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Default command runs air
CMD ["air", "-c", ".air.graphql.toml"] 