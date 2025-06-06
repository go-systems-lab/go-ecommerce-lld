# Go E-commerce Microservices

A microservices-based e-commerce system built with Go, gRPC, and GraphQL.

## 🚀 Quick Start

### Development (Hot Reload)
```bash
# Start development environment with hot reload
docker-compose -f docker-compose.dev.yml up -d

# View logs for all services (with hot reload)
docker-compose -f docker-compose.dev.yml logs -f

# View specific service logs
docker-compose -f docker-compose.dev.yml logs -f account
docker-compose -f docker-compose.dev.yml logs -f product  
docker-compose -f docker-compose.dev.yml logs -f order
```

### Production
```bash
# Start production environment
docker-compose up -d
```

## 📝 Development Features

- **Hot Reload**: Code changes in `account/`, `product/`, or `order/` automatically restart services
- **Volume Mounting**: Live code sync without rebuilding containers
- **Go Module Caching**: Faster builds with persistent module cache
- **Service Isolation**: Each service has dedicated Air config

## 🛠️ Local Development (Alternative)

```bash
# Install Air for hot reload
go install github.com/air-verse/air@latest

# Run services locally with hot reload
air -c .air.account.toml    # Account service
air -c .air.product.toml    # Product service  
air -c .air.order.toml      # Order service

# Run specific service directly
go run account/cmd/account/main.go
go run product/cmd/product/main.go
go run order/cmd/order/main.go
```

## 🌐 Services

- **GraphQL API**: `http://localhost:8080`
- **Account Service**: `localhost:5432` (PostgreSQL)
- **Product Service**: `localhost:9200` (Elasticsearch)  
- **Order Service**: `localhost:5433` (PostgreSQL)

## 📋 Development Commands

```bash
# Development with hot reload
docker-compose -f docker-compose.dev.yml up -d

# Production
docker-compose up -d

# Stop services
docker-compose down

# View logs
docker-compose logs -f [service_name]

# Rebuild specific service
docker-compose up -d --build [service_name]
```

## 🧪 Testing GraphQL

Visit `http://localhost:8080` and try:

```graphql
# Create Account
mutation {
  createAccount(account: {name: "John Doe"}) {
    id
    name
  }
}

# Create Order
mutation {
  createOrder(order: {
    accountId: "your-account-id"
    products: [
      { id: "product-id", quantity: 2 }
    ]
  }) {
    id
    totalPrice
    products {
      name
      quantity
    }
  }
}
``` 