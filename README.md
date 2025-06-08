# Go E-commerce Microservices

A production-ready microservices-based e-commerce platform built with Go, gRPC, GraphQL, and Python ML services.

## 🏗️ Architecture

- **GraphQL Gateway** - Unified API endpoint (`localhost:8080`)
- **Account Service** - User authentication & management (PostgreSQL)
- **Product Service** - Product catalog & search (Elasticsearch)
- **Order Service** - Order processing & management (PostgreSQL)
- **Recommender Service** - ML-powered product recommendations (Python/Kafka)

## 🚀 Quick Start

### Development
```bash
# Start with hot reload
docker-compose -f docker-compose.dev.yml up -d

# Access GraphQL Playground
open http://localhost:8080
```

### Production
```bash
docker-compose up -d
```

## 📋 GraphQL API

### Authentication
```graphql
# Register
mutation {
  Register(input: {
    name: "John Doe"
    email: "john@example.com"
    password: "password123"
  }) {
    token
  }
}

# Login
mutation {
  Login(input: {
    email: "john@example.com"
    password: "password123"
  }) {
    token
  }
}
```

### Products
```graphql
# Create Product
mutation {
  createProduct(product: {
    name: "iPhone 15"
    description: "Latest Apple smartphone"
    price: 999.99
  }) {
    id
    name
    price
  }
}

# Update Product
mutation {
  updateProduct(product: {
    id: "product-id"
    name: "iPhone 15 Pro"
    description: "Updated description"
    price: 1099.99
  }) {
    id
    name
    price
  }
}

# Delete Product
mutation {
  deleteProduct(id: "product-id")
}

# Get Product by ID
query {
  product(id: "product-id") {
    id
    name
    description
    price
  }
}

# Search Products
query {
  product(query: "iPhone", pagination: {skip: 0, take: 10}) {
    id
    name
    description
    price
  }
}

# Get Recommendations (based on viewed products)
query {
  product(viewedProductIds: ["product-id-1", "product-id-2"]) {
    id
    name
    price
  }
}

# Get Personalized Recommendations (for logged-in user)
query {
  product(byAccountId: true) {
    id
    name
    price
  }
}
```

### Orders
```graphql
# Create Order
mutation {
  createOrder(order: {
    products: [
      {id: "product-id", quantity: 2}
    ]
  }) {
    id
    totalPrice
    products {
      name
      quantity
      price
    }
  }
}
```

### Users
```graphql
# Get All Accounts
query {
  accounts(pagination: {skip: 0, take: 10}) {
    id
    name
    email
    orders {
      id
      totalPrice
    }
  }
}

# Get Account by ID
query {
  accounts(id: "account-id") {
    id
    name
    email
    orders {
      id
      totalPrice
    }
  }
}
```

## 🛠️ Development

### Local Services
```bash
# Generate GraphQL & protobuf
make gen && make proto

# Run individual services
go run account/cmd/account/main.go    # :50051
go run product/cmd/product/main.go    # :50052  
go run order/cmd/order/main.go        # :50053
go run graphql/main.go                # :8080
```

### Database Migrations
```bash
# Account service
make migrate-up DATABASE_URL="postgres://user:pass@localhost:5432/accounts?sslmode=disable"

# Order service  
make migrate-up DATABASE_URL="postgres://user:pass@localhost:5433/orders?sslmode=disable"
```

## 🧪 Testing

### Manual Testing
```bash
# End-to-end flow
python test_manual_flow.py

# Recommender service
python test_e2e_recommender.py
```

### Sample Data
```bash
# Debug order flow
python debug_order.py
```

## 🌐 Services & Ports

| Service | Port | Database |
|---------|------|----------|
| GraphQL Gateway | 8080 | - |
| Account Service | 50051 | PostgreSQL:5432 |
| Product Service | 50052 | Elasticsearch:9200 |
| Order Service | 50053 | PostgreSQL:5433 |
| Recommender | 50054 | Kafka, Vector DB |

## 🔧 Tech Stack

- **Backend**: Go, gRPC, GraphQL (gqlgen)
- **Database**: PostgreSQL, Elasticsearch  
- **ML**: Python, Kafka, Vector embeddings
- **DevOps**: Docker, Air (hot reload)
- **Auth**: JWT tokens

## 📚 GraphQL Schema

Visit `http://localhost:8080` for interactive schema exploration and testing.

Built with microservices best practices for scalability and maintainability. 