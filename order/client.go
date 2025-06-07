package order

import (
	"context"
	"time"

	"github.com/go-systems-lab/go-ecommerce-lld/order/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.OrderServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	C := pb.NewOrderServiceClient(conn)

	return &Client{
		conn:    conn,
		service: C,
	}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) PostOrder(ctx context.Context, accountID string, products []OrderedProduct) (*Order, error) {
	var protoProducts []*pb.OrderProduct
	for _, product := range products {
		protoProducts = append(protoProducts, &pb.OrderProduct{
			Id:       product.ID,
			Quantity: product.Quantity,
		})
	}

	r, err := c.service.PostOrder(ctx, &pb.PostOrderRequest{
		AccountId: accountID,
		Products:  protoProducts,
	})

	if err != nil {
		return nil, err
	}

	newOrder := r.Order
	newOrderCreatedAt := time.Time{}
	newOrderCreatedAt.UnmarshalBinary(newOrder.CreatedAt)

	// Convert products from gRPC response
	var responseProducts []OrderedProduct
	for _, p := range newOrder.Products {
		responseProducts = append(responseProducts, OrderedProduct{
			ID:          p.Id,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Quantity:    p.Quantity,
		})
	}

	return &Order{
		ID:         newOrder.Id,
		CreatedAt:  newOrderCreatedAt,
		AccountID:  newOrder.AccountId,
		TotalPrice: newOrder.TotalPrice,
		Products:   responseProducts,
	}, nil
}

func (c *Client) GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error) {
	r, err := c.service.GetOrdersForAccount(ctx, &pb.GetOrdersForAccountRequest{
		AccountId: accountID,
	})
	if err != nil {
		return nil, err
	}

	// Create response orders
	orders := []Order{}
	for _, orderProto := range r.Orders {
		newOrder := Order{
			ID:         orderProto.Id,
			TotalPrice: orderProto.TotalPrice,
			AccountID:  orderProto.AccountId,
		}
		newOrder.CreatedAt = time.Time{}
		newOrder.CreatedAt.UnmarshalBinary(orderProto.CreatedAt)

		var products []OrderedProduct
		for _, p := range orderProto.Products {
			products = append(products, OrderedProduct{
				ID:          p.Id,
				Quantity:    p.Quantity,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
			})
		}
		newOrder.Products = products

		orders = append(orders, newOrder)
	}
	return orders, nil
}
