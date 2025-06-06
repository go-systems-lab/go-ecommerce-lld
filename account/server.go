package account

import (
	"context"
	"fmt"
	"net"

	"github.com/go-systems-lab/go-ecommerce-lld/account/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	pb.UnimplementedAccountServiceServer
	service Service
}

func ListenGRPC(s Service, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	srv := grpc.NewServer()
	pb.RegisterAccountServiceServer(srv, &grpcServer{service: s, UnimplementedAccountServiceServer: pb.UnimplementedAccountServiceServer{}})
	reflection.Register(srv)
	return srv.Serve(lis)
}

func (s *grpcServer) RegisterAccount(ctx context.Context, req *pb.RegisterRequest) (*pb.AuthResponse, error) {
	token, err := s.service.Register(ctx, req.Name, req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	return &pb.AuthResponse{Token: token}, nil
}

func (s *grpcServer) LoginAccount(ctx context.Context, req *pb.LoginRequest) (*pb.AuthResponse, error) {
	token, err := s.service.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	return &pb.AuthResponse{Token: token}, nil
}

func (s *grpcServer) GetAccount(ctx context.Context, req *pb.GetAccountRequest) (*pb.AccountResponse, error) {
	a, err := s.service.GetAccountByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.AccountResponse{Account: &pb.Account{
		Id:   a.ID,
		Name: a.Name,
	}}, nil
}

func (s *grpcServer) GetAccounts(ctx context.Context, req *pb.GetAccountsRequest) (*pb.GetAccountsResponse, error) {
	a, err := s.service.ListAccounts(ctx, req.Skip, req.Take)
	if err != nil {
		return nil, err
	}

	var accounts []*pb.Account
	for _, a := range a {
		accounts = append(accounts, &pb.Account{
			Id:   a.ID,
			Name: a.Name,
		})
	}
	return &pb.GetAccountsResponse{Accounts: accounts}, nil
}
