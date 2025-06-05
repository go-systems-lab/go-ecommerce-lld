package product

import (
	"context"

	"github.com/google/uuid"
)

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type Service interface {
	PostProduct(ctx context.Context, name, description string, price float64) (*Product, error)
	GetProduct(ctx context.Context, id string) (*Product, error)
	GetProducts(ctx context.Context, skip, take uint64) ([]Product, error)
	GetProductsWithIDs(ctx context.Context, ids []string) ([]Product, error)
	SearchProducts(ctx context.Context, query string, skip, take uint64) ([]Product, error)
	UpdateProduct(ctx context.Context, id, name, description string, price float64) (*Product, error)
	DeleteProduct(ctx context.Context, productId string) error
}

type productService struct {
	repo Repository
}

func NewProductService(repo Repository) Service {
	return &productService{repo: repo}
}

func (p productService) PostProduct(ctx context.Context, name, description string, price float64) (*Product, error) {
	product := Product{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		Price:       price,
	}

	if err := p.repo.PutProduct(ctx, product); err != nil {
		return nil, err
	}

	return &product, nil
}

func (p productService) GetProduct(ctx context.Context, id string) (*Product, error) {
	return p.repo.GetProductById(ctx, id)
}

func (p productService) GetProducts(ctx context.Context, skip, take uint64) ([]Product, error) {
	return p.repo.ListProducts(ctx, skip, take)
}

func (p productService) GetProductsWithIDs(ctx context.Context, ids []string) ([]Product, error) {
	return p.repo.ListProductsWithIds(ctx, ids)
}

func (p productService) SearchProducts(ctx context.Context, query string, skip, take uint64) ([]Product, error) {
	return p.repo.SearchProducts(ctx, query, skip, take)
}

func (p productService) UpdateProduct(ctx context.Context, id, name, description string, price float64) (*Product, error) {
	updatedProduct := Product{
		ID:          id,
		Name:        name,
		Description: description,
		Price:       price,
	}

	if err := p.repo.UpdateProduct(ctx, updatedProduct); err != nil {
		return nil, err
	}

	return &updatedProduct, nil
}

func (p productService) DeleteProduct(ctx context.Context, productId string) error {
	return p.repo.DeleteProduct(ctx, productId)
}
