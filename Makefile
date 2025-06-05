gen:
	go run github.com/99designs/gqlgen generate

dep:
	go mod tidy && go fmt

run:
	go run graphql/*.go

proto:
	protoc --go_out=./account/pb --go_opt=paths=source_relative --go-grpc_out=./account/pb --go-grpc_opt=paths=source_relative --proto_path=./account account.proto
	protoc --go_out=./product/pb --go_opt=paths=source_relative --go-grpc_out=./product/pb --go-grpc_opt=paths=source_relative --proto_path=./product product.proto

# Docker targets - Database
db-build:
	docker build -f account/db.dockerfile -t ecommerce-account-db ./account

db-run:
	docker run -d --name ecommerce-account-db -p 5432:5432 ecommerce-account-db

db-stop:
	docker stop ecommerce-account-db || true
	docker rm ecommerce-account-db || true

# Docker targets - Account Service
app-build:
	docker build -f account/app.dockerfile -t ecommerce-account-app .

app-run:
	docker run -d --name ecommerce-account-app -p 50051:50051 \
		-e DATABASE_URL="postgres://postgres:postgres@host.docker.internal:5432/ecommerce_account?sslmode=disable" \
		-e PORT=50051 \
		ecommerce-account-app

app-stop:
	docker stop ecommerce-account-app || true
	docker rm ecommerce-account-app || true

# Full Docker Stack
docker-up:
	@echo "🚀 Starting full Docker stack..."
	@make db-stop app-stop || true
	@make db-build app-build
	@make db-run
	@echo "⏳ Waiting for database to be ready..."
	@sleep 8
	@make migrate-up DATABASE_URL="postgres://postgres:postgres@localhost:5432/ecommerce_account?sslmode=disable"
	@make app-run
	@echo "⏳ Waiting for account service to be ready..."
	@sleep 3
	@echo "✅ Full stack is ready!"
	@echo "📊 Database: localhost:5432"
	@echo "🔧 Account Service (gRPC): localhost:50051"
	@echo "🌐 Now run: make run (for GraphQL server on :8080)"

docker-down:
	@echo "🛑 Stopping full Docker stack..."
	@make app-stop
	@make db-stop
	@echo "✅ Stack stopped!"

docker-logs-db:
	docker logs -f ecommerce-account-db

docker-logs-app:
	docker logs -f ecommerce-account-app

# Account service targets
account-run:
	go run account/cmd/account/main.go

# Full stack testing targets
test-setup:
	@echo "Setting up full testing environment..."
	@make db-stop || true
	@make db-build
	@make db-run
	@echo "Waiting for database to be ready..."
	@sleep 5
	@make migrate-up DATABASE_URL="postgres://postgres:postgres@localhost:5432/ecommerce_account?sslmode=disable"
	@echo "✅ Test environment ready!"
	@echo "Now run: make account-run (in terminal 1) and make run (in terminal 2)"

test-teardown:
	@echo "Tearing down test environment..."
	@make db-stop
	@echo "✅ Test environment cleaned up!"

# Complete testing workflow
test-full:
	@echo "🧪 Starting complete testing workflow..."
	@make docker-up
	@echo ""
	@echo "🎯 Testing GraphQL API Gateway:"
	@echo "1. Run: make run (to start GraphQL server)"
	@echo "2. Open: http://localhost:8080 (GraphQL Playground)"
	@echo "3. Try sample queries (see test-queries target)"
	@echo ""
	@make test-queries

test-queries:
	@echo "📝 Sample GraphQL Queries to test:"
	@echo ""
	@echo "🔸 Create Account Mutation:"
	@echo 'mutation {'
	@echo '  createAccount(account: {name: "John Doe"}) {'
	@echo '    id'
	@echo '    name'
	@echo '  }'
	@echo '}'
	@echo ""
	@echo "🔸 Get All Accounts Query:"
	@echo 'query {'
	@echo '  accounts(pagination: {skip: 0, take: 10}) {'
	@echo '    id'
	@echo '    name'
	@echo '  }'
	@echo '}'
	@echo ""
	@echo "🔸 Get Account by ID Query:"
	@echo 'query {'
	@echo '  accounts(id: "your-account-id-here") {'
	@echo '    id'
	@echo '    name'
	@echo '  }'
	@echo '}'

# Migration targets for account service
migrate-create:
	@if [ -z "$(name)" ]; then echo "Usage: make migrate-create name=migration_name"; exit 1; fi
	migrate create -ext sql -dir account/migrations -seq $(name)

migrate-up:
	@if [ -z "$(DATABASE_URL)" ]; then echo "DATABASE_URL is required. Example: make migrate-up DATABASE_URL='postgres://user:password@localhost:5432/dbname?sslmode=disable'"; exit 1; fi
	migrate -path account/migrations -database "$(DATABASE_URL)" up

migrate-down:
	@if [ -z "$(DATABASE_URL)" ]; then echo "DATABASE_URL is required. Example: make migrate-down DATABASE_URL='postgres://user:password@localhost:5432/dbname?sslmode=disable'"; exit 1; fi
	migrate -path account/migrations -database "$(DATABASE_URL)" down

migrate-up-by:
	@if [ -z "$(DATABASE_URL)" ]; then echo "DATABASE_URL is required"; exit 1; fi
	@if [ -z "$(steps)" ]; then echo "Usage: make migrate-up-by steps=N DATABASE_URL='...'"; exit 1; fi
	migrate -path account/migrations -database "$(DATABASE_URL)" up $(steps)

migrate-down-by:
	@if [ -z "$(DATABASE_URL)" ]; then echo "DATABASE_URL is required"; exit 1; fi
	@if [ -z "$(steps)" ]; then echo "Usage: make migrate-down-by steps=N DATABASE_URL='...'"; exit 1; fi
	migrate -path account/migrations -database "$(DATABASE_URL)" down $(steps)

migrate-force:
	@if [ -z "$(DATABASE_URL)" ]; then echo "DATABASE_URL is required"; exit 1; fi
	@if [ -z "$(version)" ]; then echo "Usage: make migrate-force version=VERSION DATABASE_URL='...'"; exit 1; fi
	migrate -path account/migrations -database "$(DATABASE_URL)" force $(version)

migrate-version:
	@if [ -z "$(DATABASE_URL)" ]; then echo "DATABASE_URL is required"; exit 1; fi
	migrate -path account/migrations -database "$(DATABASE_URL)" version

migrate-drop:
	@if [ -z "$(DATABASE_URL)" ]; then echo "DATABASE_URL is required"; exit 1; fi
	@echo "WARNING: This will drop all tables and data!"
	@read -p "Are you sure? (y/N): " confirm && [ "$$confirm" = "y" ]
	migrate -path account/migrations -database "$(DATABASE_URL)" drop -f

migrate-help:
	@echo "Migration Commands:"
	@echo "  migrate-create name=<migration_name>    - Create a new migration"
	@echo "  migrate-up DATABASE_URL=<url>           - Run all pending migrations"
	@echo "  migrate-down DATABASE_URL=<url>         - Rollback the last migration"
	@echo "  migrate-up-by steps=N DATABASE_URL=<url> - Run N migrations up"
	@echo "  migrate-down-by steps=N DATABASE_URL=<url> - Run N migrations down"
	@echo "  migrate-force version=V DATABASE_URL=<url> - Force set migration version"
	@echo "  migrate-version DATABASE_URL=<url>      - Show current migration version"
	@echo "  migrate-drop DATABASE_URL=<url>         - Drop all tables (DANGEROUS!)"
	@echo ""
	@echo "Example DATABASE_URL: postgres://user:password@localhost:5432/dbname?sslmode=disable"