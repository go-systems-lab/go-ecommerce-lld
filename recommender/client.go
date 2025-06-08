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

func (c *Client) GetRecommendationForUserId(ctx context.Context, userID string, skip uint64, take uint64) (*pb.RecommendationResponse, error) {
	return c.service.GetRecommendationsForUserId(
		ctx,
		&pb.RecommendationRequestForUserId{
			UserId: userID,
			Skip:   skip,
			Take:   take,
		},
	)
}

func (c *Client) GetRecommendationOnViews(ctx context.Context, ids []string, skip uint64, take uint64) (*pb.RecommendationResponse, error) {
	return c.service.GetRecommendationsOnViews(
		ctx,
		&pb.RecommendationRequestOnViews{
			Ids:  ids,
			Skip: skip,
			Take: take,
		},
	)
}
