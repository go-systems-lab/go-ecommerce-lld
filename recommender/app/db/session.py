"""
Database session configuration for the recommender service.

This module provides database engine and session management using SQLAlchemy.
"""

from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker, Session
from typing import Generator
from config.settings import DATABASE_URL

from .models import Base

# Create engine with connection pooling
engine = create_engine(
    DATABASE_URL,
    echo=False,  # Set to True for SQL query logging in development
    pool_pre_ping=True,
    pool_recycle=300,
)

# Create sessionmaker
SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)


def create_tables():
    """Create all database tables."""
    Base.metadata.create_all(bind=engine)


def get_db() -> Generator[Session, None, None]:
    """
    Dependency function to get database session.
    
    Yields:
        Session: SQLAlchemy database session
    """
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()


def get_db_session() -> Session:
    """
    Get a database session (for direct use).
    
    Returns:
        Session: SQLAlchemy database session
    """
    return SessionLocal()
