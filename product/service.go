package product

import (
	"context"
	"errors"
	"log"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	AccountID   string  `json:"accountId"`
}

type Service interface {
	PostProduct(ctx context.Context, name, description string, price float64, accountId string) (*Product, error)
	GetProduct(ctx context.Context, id string) (*Product, error)
	GetProducts(ctx context.Context, skip, take uint64) ([]Product, error)
	GetProductsWithIDs(ctx context.Context, ids []string) ([]Product, error)
	SearchProducts(ctx context.Context, query string, skip, take uint64) ([]Product, error)
	UpdateProduct(ctx context.Context, id, name, description string, price float64, accountId string) (*Product, error)
	DeleteProduct(ctx context.Context, productId string, accountId string) error
}

type productService struct {
	repo     Repository
	producer sarama.AsyncProducer
}

func NewProductService(repo Repository, producer sarama.AsyncProducer) Service {
	return &productService{repo: repo, producer: producer}
}

func (p productService) PostProduct(ctx context.Context, name, description string, price float64, accountId string) (*Product, error) {
	log.Printf("PostProduct called with accountId: %s (type: %T)", accountId, accountId)

	product := Product{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		Price:       price,
		AccountID:   accountId,
	}

	log.Printf("Created product struct: %+v", product)

	if err := p.repo.PutProduct(ctx, product); err != nil {
		log.Printf("Error from repository.PutProduct: %v", err)
		return nil, err
	}

	go func() {
		err := p.SendMessageToRecommender(Event{
			Type: "product_created",
			Data: EventData{
				ID:          &product.ID,
				Name:        &product.Name,
				Description: &product.Description,
				Price:       &product.Price,
				AccountID:   &accountId,
			},
		}, "product_events")
		if err != nil {
			log.Printf("Error sending message to recommender: %v", err)
		}
	}()

	log.Printf("Successfully stored product in repository")
	return &product, nil
}

func (p productService) GetProduct(ctx context.Context, id string) (*Product, error) {
	product, err := p.repo.GetProductById(ctx, id)
	if err != nil {
		return nil, err
	}

	go func() {
		err := p.SendMessageToRecommender(Event{
			Type: "product_retrieved",
			Data: EventData{
				ID:        &product.ID,
				AccountID: &product.AccountID,
			},
		}, "interaction_events")
		if err != nil {
			log.Printf("Error sending message to recommender: %v", err)
		}
	}()

	return product, nil
}

func (p productService) GetProducts(ctx context.Context, skip, take uint64) ([]Product, error) {
	products, err := p.repo.ListProducts(ctx, skip, take)
	if err != nil {
		return nil, err
	}

	log.Printf("GetProducts: Retrieved %d products from repository", len(products))
	for i, product := range products {
		log.Printf("Product %d: ID=%s, Name=%s, AccountID=%s", i, product.ID, product.Name, product.AccountID)
	}

	return products, nil
}

func (p productService) GetProductsWithIDs(ctx context.Context, ids []string) ([]Product, error) {
	return p.repo.ListProductsWithIds(ctx, ids)
}

func (p productService) SearchProducts(ctx context.Context, query string, skip, take uint64) ([]Product, error) {
	return p.repo.SearchProducts(ctx, query, skip, take)
}

func (p productService) UpdateProduct(ctx context.Context, id, name, description string, price float64, accountId string) (*Product, error) {
	product, err := p.repo.GetProductById(ctx, id)
	if err != nil {
		return nil, err
	}

	if product.AccountID != accountId {
		return nil, errors.New("unauthorized")
	}

	updatedProduct := Product{
		ID:          id,
		Name:        name,
		Description: description,
		Price:       price,
		AccountID:   accountId,
	}

	if err = p.repo.UpdateProduct(ctx, updatedProduct); err != nil {
		return nil, err
	}

	go func() {
		err := p.SendMessageToRecommender(Event{
			Type: "product_updated",
			Data: EventData{
				ID:          &updatedProduct.ID,
				Name:        &updatedProduct.Name,
				Description: &updatedProduct.Description,
				Price:       &updatedProduct.Price,
				AccountID:   &accountId,
			},
		}, "product_events")
		if err != nil {
			log.Printf("Error sending message to recommender: %v", err)
		}
	}()

	return &updatedProduct, nil
}

func (p productService) DeleteProduct(ctx context.Context, productId string, accountId string) error {
	product, err := p.repo.GetProductById(ctx, productId)
	if err != nil {
		return err
	}

	if product.AccountID != accountId {
		return errors.New("unauthorized")
	}

	go func() {

		err = p.SendMessageToRecommender(Event{
			Type: "product_deleted",
			Data: EventData{
				ID: &product.ID,
			},
		}, "product_events")
		if err != nil {
			log.Printf("Error sending message to recommender: %v", err)
		}
	}()

	if err = p.repo.DeleteProduct(ctx, productId); err != nil {
		return err
	}

	return nil
}
