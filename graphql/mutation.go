package main

import (
	"context"
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-systems-lab/go-ecommerce-lld/account"
	"github.com/go-systems-lab/go-ecommerce-lld/order"
)

var ErrInvalidParameter = errors.New("invalid parameter")

type mutationResolver struct {
	server *Server
}

func (r *mutationResolver) Register(ctx context.Context, input RegisterInput) (*AuthResponse, error) {
	// Create account via microservice
	token, err := r.server.accountClient.Register(ctx, input.Name, input.Email, input.Password)
	if err != nil {
		return nil, err
	}

	ginContext, ok := ctx.Value("GinContextKey").(*gin.Context)
	if !ok {
		return nil, errors.New("gin context not found")
	}

	ginContext.SetCookie("token", token, 3600, "/", "localhost", false, true)
	return &AuthResponse{Token: token}, nil
}

func (r *mutationResolver) Login(ctx context.Context, input LoginInput) (*AuthResponse, error) {
	token, err := r.server.accountClient.Login(ctx, input.Email, input.Password)
	if err != nil {
		return nil, err
	}

	ginContext, ok := ctx.Value("GinContextKey").(*gin.Context)
	if !ok {
		return nil, errors.New("gin context not found")
	}

	ginContext.SetCookie("token", token, 3600, "/", "localhost", false, true)
	return &AuthResponse{Token: token}, nil
}

func (r *mutationResolver) CreateProduct(ctx context.Context, product CreateProductInput) (*Product, error) {
	log.Println("CreateProduct called with:", product)

	accountId := account.GetUserId(ctx)
	log.Println("Got accountId:", accountId)
	if accountId == "" {
		return nil, errors.New("unauthorized")
	}

	log.Println("Calling productClient.PostProduct with accountId:", accountId)
	createdProduct, err := r.server.productClient.PostProduct(ctx, product.Name, product.Description, product.Price, accountId)
	if err != nil {
		log.Println("Error from productClient.PostProduct:", err)
		return nil, err
	}

	log.Println("Successfully created product:", createdProduct)
	return &Product{
		ID:          createdProduct.ID,
		Name:        createdProduct.Name,
		Description: createdProduct.Description,
		Price:       createdProduct.Price,
		AccountID:   accountId,
	}, nil
}

func (r *mutationResolver) CreateOrder(ctx context.Context, in OrderInput) (*Order, error) {
	log.Printf("CreateOrder called with %d products", len(in.Products))

	var products []order.OrderedProduct
	for _, p := range in.Products {
		if p.Quantity <= 0 {
			return nil, ErrInvalidParameter
		}
		log.Printf("Adding product to order: ID=%s, Quantity=%d", p.ID, p.Quantity)
		products = append(products, order.OrderedProduct{
			ID:       p.ID,
			Quantity: uint32(p.Quantity),
		})
	}

	accountId := account.GetUserId(ctx)
	if accountId == "" {
		return nil, errors.New("unauthorized")
	}

	log.Printf("Calling orderClient.PostOrder with accountId=%s and %d products", accountId, len(products))
	o, err := r.server.orderClient.PostOrder(ctx, accountId, products)
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

func (r *mutationResolver) UpdateProduct(ctx context.Context, product UpdateProductInput) (*Product, error) {
	accountId := account.GetUserId(ctx)
	if accountId == "" {
		return nil, errors.New("unauthorized")
	}

	updatedProduct, err := r.server.productClient.UpdateProduct(ctx, product.ID, product.Name, product.Description, product.Price, accountId)
	if err != nil {
		return nil, err
	}

	return &Product{
		ID:          updatedProduct.ID,
		Name:        updatedProduct.Name,
		Description: updatedProduct.Description,
		Price:       updatedProduct.Price,
		AccountID:   accountId,
	}, nil
}

func (r *mutationResolver) DeleteProduct(ctx context.Context, id string) (*bool, error) {
	accountId := account.GetUserId(ctx)
	if accountId == "" {
		result := false
		return &result, errors.New("unauthorized")
	}

	err := r.server.productClient.DeleteProduct(ctx, id, accountId)
	if err != nil {
		result := false
		return &result, err
	}

	result := true
	return &result, nil
}
