package main

import (
	"context"
	pb "gRPC-Project/proto/gen"
)

type Service struct {
	pb.UnimplementedAuthServiceServer
}

func (s *Service) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {

	return &pb.RegisterResponse{}, nil
}

func main() {

}
