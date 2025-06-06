package product

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/go-systems-lab/go-ecommerce-lld/product/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
)

type grpcServer struct {
	pb.UnimplementedProductServiceServer
	service Service
}

func ListenGRPC(s Service, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	srv := grpc.NewServer()
	pb.RegisterProductServiceServer(srv, &grpcServer{service: s, UnimplementedProductServiceServer: pb.UnimplementedProductServiceServer{}})
	reflection.Register(srv)
	return srv.Serve(lis)
}

func (s *grpcServer) PostProduct(ctx context.Context, r *pb.CreateProductRequest) (*pb.ProductResponse, error) {
	p, err := s.service.PostProduct(ctx, r.GetName(), r.GetDescription(), r.GetPrice(), r.GetAccountId())
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &pb.ProductResponse{Product: &pb.Product{
		Id:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		AccountId:   p.AccountID,
	}}, nil
}

func (s *grpcServer) GetProduct(ctx context.Context, req *pb.ProductByIdRequest) (*pb.ProductResponse, error) {
	p, err := s.service.GetProduct(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.ProductResponse{Product: &pb.Product{
		Id:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		AccountId:   p.AccountID,
	}}, nil
}

func (s *grpcServer) GetProducts(ctx context.Context, req *pb.GetProductsRequest) (*pb.ProductsResponse, error) {
	var products []Product

	if len(req.Ids) > 0 {
		res, err := s.service.GetProductsWithIDs(ctx, req.Ids)
		if err != nil {
			return nil, err
		}
		products = res
	} else if req.Query != "" {
		res, err := s.service.SearchProducts(ctx, req.Query, req.Skip, req.Take)
		if err != nil {
			return nil, err
		}
		products = res
	} else {
		res, err := s.service.GetProducts(ctx, req.Skip, req.Take)
		if err != nil {
			return nil, err
		}
		products = res
	}

	var pbProducts []*pb.Product
	for _, p := range products {
		pbProducts = append(pbProducts, &pb.Product{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			AccountId:   p.AccountID,
		})
	}

	return &pb.ProductsResponse{Products: pbProducts}, nil
}

func (s *grpcServer) UpdateProduct(ctx context.Context, r *pb.UpdateProductRequest) (*pb.ProductResponse, error) {
	p, err := s.service.UpdateProduct(ctx, r.GetId(), r.GetName(), r.GetDescription(), r.GetPrice(), r.GetAccountId())
	if err != nil {
		return nil, err
	}

	return &pb.ProductResponse{Product: &pb.Product{
		Id:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		AccountId:   p.AccountID,
	}}, nil
}

func (s *grpcServer) DeleteProduct(ctx context.Context, r *pb.DeleteProductRequest) (*emptypb.Empty, error) {
	err := s.service.DeleteProduct(ctx, r.GetProductId(), r.GetAccountId())
	return &emptypb.Empty{}, err
}
