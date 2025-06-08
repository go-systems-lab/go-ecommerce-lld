package recommender

import (
	"context"

	"github.com/go-systems-lab/go-ecommerce-lld/recommender/generated/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.RecommenderServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:    conn,
		service: pb.NewRecommenderServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) GetRecommendation(ctx context.Context, userID string) (*pb.RecommendationResponse, error) {
	return c.service.GetRecommendations(
		ctx,
		&pb.RecommendationRequest{
			UserId: userID,
		},
	)
}
