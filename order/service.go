package order

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID           string
	CreatedAt    time.Time
	AccountID    string
	TotalPrice   float64
	Products     []OrderedProduct
	productInfos []ProductsInfo
}

type ProductsInfo struct {
	ID        string
	OrderID   string
	ProductID string
	Quantity  int
}

type OrderedProduct struct {
	ID          string
	Name        string
	Description string
	Price       float64
	Quantity    uint32
}

type Service interface {
	PostOrder(ctx context.Context, accountID string, totalPrice float64, products []OrderedProduct) (*Order, error)
	GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error)
}

type orderService struct {
	repository Repository
}

func NewOrderService(repository Repository) Service {
	return &orderService{repository: repository}
}

func (o orderService) PostOrder(ctx context.Context, accountID string, totalPrice float64, products []OrderedProduct) (*Order, error) {
	order := Order{
		ID:         uuid.New().String(),
		CreatedAt:  time.Now().UTC(),
		AccountID:  accountID,
		TotalPrice: totalPrice,
		Products:   products,
	}

	err := o.repository.PutOrder(ctx, order)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (o orderService) GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error) {
	return o.repository.GetOrdersForAccount(ctx, accountID)
}
