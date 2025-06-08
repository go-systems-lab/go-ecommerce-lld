# Recommender Service Database Setup

## Overview
Clean and professional database setup using SQLAlchemy and Alembic for the recommender service.

## Database Structure
- **Products**: Core product information synchronized from product service
- **Interactions**: User behavior tracking (views, purchases, ratings, etc.)

## Quick Start

### 1. Start Database
```bash
# Start recommender database service
docker-compose -f docker-compose.dev.yml up recommender_db -d
```

### 2. Run Migrations
```bash
# Activate virtual environment
source venv/bin/activate

# Run database migrations
alembic upgrade head
```

### 3. Test Setup
```bash
# Test database connectivity and operations
python test_db.py
```

## Database Configuration
- **Host**: localhost:5434 (development)
- **Database**: recommender_db
- **User**: postgres
- **Password**: postgres

## Environment Variables
```bash
DATABASE_URL=postgresql://postgres:postgres@localhost:5434/recommender_db
```

## Migration Commands
```bash
# Create new migration
alembic revision --autogenerate -m "Description"

# Apply migrations
alembic upgrade head

# Rollback migration
alembic downgrade -1
```

## Models
- `Product`: Product information with relationships
- `Interaction`: User behavior tracking
- `InteractionType`: Constants for interaction types 