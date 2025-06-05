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
	// TODO: Implement
	return nil, nil
}
