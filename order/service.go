package order

import (
	"context"
	"log"
	"time"

	"github.com/IBM/sarama"
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
	producer   sarama.AsyncProducer
}

func NewOrderService(repository Repository, producer sarama.AsyncProducer) Service {
	return &orderService{repository: repository, producer: producer}
}

func (o orderService) PostOrder(ctx context.Context, accountID string, totalPrice float64, products []OrderedProduct) (*Order, error) {
	log.Printf("PostOrder called with accountID: %s, %d products", accountID, len(products))

	order := Order{
		ID:         uuid.New().String(),
		CreatedAt:  time.Now().UTC(),
		AccountID:  accountID,
		TotalPrice: totalPrice,
		Products:   products,
	}

	log.Printf("Created order with ID: %s", order.ID)

	err := o.repository.PutOrder(ctx, order)
	if err != nil {
		log.Printf("Error storing order in repository: %v", err)
		return nil, err
	}

	log.Printf("Order stored in repository successfully, sending %d interaction events", len(products))

	// send to recommendation service
	go func() {
		for i, product := range products {
			log.Printf("Sending interaction event %d/%d for product %s", i+1, len(products), product.ID)

			event := Event{
				Type: "purchase",
				EventData: EventData{
					AccountId: accountID,
					ProductId: product.ID,
				},
			}

			log.Printf("Sending event: %+v", event)

			err = o.SendMessageToRecommender(event, "interaction_events")
			if err != nil {
				log.Printf("Error sending message to recommender: %v", err)
			} else {
				log.Printf("Successfully sent interaction event for product %s", product.ID)
			}
		}
		log.Printf("Finished sending all %d interaction events", len(products))
	}()

	return &order, nil
}

func (o orderService) GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error) {
	return o.repository.GetOrdersForAccount(ctx, accountID)
}
