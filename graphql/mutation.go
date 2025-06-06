package main

import (
	"context"
	"errors"
	"log"

	"github.com/go-systems-lab/go-ecommerce-lld/order"
)

var ErrInvalidParameter = errors.New("invalid parameter")

type mutationResolver struct {
	server *Server
}

func (r *mutationResolver) CreateAccount(ctx context.Context, account AccountInput) (*Account, error) {
	// Create account via microservice
	createdAccount, err := r.server.accountClient.PostAccount(ctx, account.Name)
	if err != nil {
		return nil, err
	}

	return &Account{
		ID:   createdAccount.ID,
		Name: createdAccount.Name,
	}, nil
}

func (r *mutationResolver) CreateProduct(ctx context.Context, product ProductInput) (*Product, error) {
	createdProduct, err := r.server.productClient.PostProduct(ctx, product.Name, product.Description, product.Price)
	if err != nil {
		return nil, err
	}

	return &Product{
		ID:          createdProduct.ID,
		Name:        createdProduct.Name,
		Description: createdProduct.Description,
		Price:       createdProduct.Price,
	}, nil
}

func (r *mutationResolver) CreateOrder(ctx context.Context, in OrderInput) (*Order, error) {
	var products []order.OrderedProduct
	for _, p := range in.Products {
		if p.Quantity <= 0 {
			return nil, ErrInvalidParameter
		}
		products = append(products, order.OrderedProduct{
			ID:       p.ID,
			Quantity: uint32(p.Quantity),
		})
	}
	o, err := r.server.orderClient.PostOrder(ctx, in.AccountID, products)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Convert products for GraphQL response
	var orderProducts []*OrderedProduct
	for _, p := range o.Products {
		orderProducts = append(orderProducts, &OrderedProduct{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Quantity:    int(p.Quantity),
		})
	}

	return &Order{
		ID:         o.ID,
		CreatedAt:  o.CreatedAt,
		TotalPrice: o.TotalPrice,
		Products:   orderProducts,
	}, nil
}
