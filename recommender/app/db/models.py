"""
Database models for the recommender service.

This module contains SQLAlchemy models for products and user interactions
used by the recommendation engine.
"""

from datetime import datetime
from typing import List

from sqlalchemy import String, Float, Integer, DateTime, func, ForeignKey
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import Mapped, mapped_column, relationship

Base = declarative_base()


class Product(Base):
    """
    Product model representing products in the recommender system.
    
    This model stores product information that can be synchronized from
    the main product service and used for generating recommendations.
    """
    __tablename__ = "products"

    id: Mapped[str] = mapped_column(String, primary_key=True)
    name: Mapped[str] = mapped_column(String, nullable=False)
    description: Mapped[str] = mapped_column(String, nullable=False)
    price: Mapped[float] = mapped_column(Float, nullable=False)
    account_id: Mapped[str] = mapped_column(String, nullable=False)  # Changed to String for consistency
    created_at: Mapped[datetime] = mapped_column(DateTime, default=func.now())
    updated_at: Mapped[datetime] = mapped_column(DateTime, default=func.now(), onupdate=func.now())

    # Relationships
    interactions: Mapped[List["Interaction"]] = relationship(
        "Interaction", 
        back_populates="product", 
        cascade="all, delete-orphan"
    )

    def __repr__(self) -> str:
        return f"<Product(id='{self.id}', name='{self.name}', price={self.price})>"


class Interaction(Base):
    """
    User interaction model for tracking user behavior with products.
    
    This model stores various types of user interactions (views, purchases, ratings)
    that are used to generate personalized recommendations.
    """
    __tablename__ = "interactions"

    id: Mapped[int] = mapped_column(Integer, primary_key=True, autoincrement=True)
    user_id: Mapped[str] = mapped_column(String, nullable=False)
    product_id: Mapped[str] = mapped_column(String, ForeignKey("products.id"), nullable=False)
    interaction_type: Mapped[str] = mapped_column(String, nullable=False)  # view, purchase, rating, etc.
    rating: Mapped[float] = mapped_column(Float, nullable=True)  # Optional rating value
    timestamp: Mapped[datetime] = mapped_column(DateTime, default=func.now())

    # Relationships
    product: Mapped["Product"] = relationship("Product", back_populates="interactions")

    def __repr__(self) -> str:
        return f"<Interaction(id={self.id}, user='{self.user_id}', product='{self.product_id}', type='{self.interaction_type}')>"


# Interaction types enum-like constants
class InteractionType:
    """Constants for different types of user interactions."""
    VIEW = "view"
    PURCHASE = "purchase"
    RATING = "rating"
    ADD_TO_CART = "add_to_cart"
    REMOVE_FROM_CART = "remove_from_cart"
    WISHLIST = "wishlist" 