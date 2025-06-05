package main

import (
	"context"
)

type queryResolver struct {
	server *Server
}

func (r *queryResolver) Accounts(ctx context.Context, pagination *PaginationInput, id *string) ([]*Account, error) {
	// If specific ID is requested, get single account
	if id != nil {
		account, err := r.server.accountClient.GetAccount(ctx, *id)
		if err != nil {
			return nil, err
		}
		return []*Account{
			{
				ID:   account.ID,
				Name: account.Name,
			},
		}, nil
	}

	// Otherwise get paginated list
	skip := uint64(0)
	take := uint64(10) // default

	if pagination != nil {
		skip = uint64(pagination.Skip)
		take = uint64(pagination.Take)
	}

	accounts, err := r.server.accountClient.GetAccounts(ctx, skip, take)
	if err != nil {
		return nil, err
	}

	var result []*Account
	for _, acc := range accounts {
		result = append(result, &Account{
			ID:   acc.ID,
			Name: acc.Name,
		})
	}

	return result, nil
}

func (r *queryResolver) Product(ctx context.Context, pagination *PaginationInput, query, id *string) ([]*Product, error) {
	if id != nil {
		product, err := r.server.productClient.GetProduct(ctx, *id)
		if err != nil {
			return nil, err
		}
		return []*Product{
			{
				ID:          product.ID,
				Name:        product.Name,
				Description: product.Description,
				Price:       product.Price,
			},
		}, nil
	}
	skip, take := uint64(0), uint64(10)
	if pagination != nil {
		skip, take = pagination.bounds()
	}

	q := ""
	if query != nil {
		q = *query
	}

	products, err := r.server.productClient.GetProducts(ctx, skip, take, nil, q)
	if err != nil {
		return nil, err
	}

	var result []*Product
	for _, p := range products {
		result = append(result, &Product{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		})
	}
	return result, nil
}

func (p PaginationInput) bounds() (uint64, uint64) {
	skipValue := uint64(0)
	takeValue := uint64(100)
	if p.Skip != 0 {
		skipValue = uint64(p.Skip)
	}
	if p.Take != 0 {
		takeValue = uint64(p.Take)
	}
	return skipValue, takeValue
}
