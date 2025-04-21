package grpc

import (
	"context"
	"fmt"
	pb "gRPC-Project/internal/grpc/gen"
	"gRPC-Project/internal/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type authServer struct {
	pb.UnimplementedAuthServiceServer
	service service.AuthService
}

func NewAuthServer(service service.AuthService) *authServer {
	return &authServer{
		service: service,
	}
}

func (s *authServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	// Проверяем входные данные

	if req.Username == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "username and password are required")
	}

	// Вызываем сервисный метод
	userID, err := s.service.Register(ctx, req.Username, req.Password)
	if err != nil {
		// Возвращаем соответствующую ошибку gRPC
		if err.Error() == "user already exists" {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Error(codes.Internal, "failed to register user")
	}

	// Возвращаем успешный ответ
	return &pb.RegisterResponse{
		UserId:  userID,
		Message: fmt.Sprintf("%s:%s", "User registered successfully", userID),
	}, nil
}
func (s *authServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	// Проверяем входные данные
	if req.Username == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "username and password are required")
	}

	// Вызываем сервисный метод
	token, userID, err := s.service.Login(ctx, req.Username, req.Password)
	if err != nil {
		// Возвращаем соответствующую ошибку gRPC
		return nil, status.Error(codes.Unauthenticated, "invalid username or password")
	}

	// Возвращаем успешный ответ с токеном
	return &pb.LoginResponse{
		Token:  token,
		UserId: userID,
	}, nil
}

func (s *authServer) CheckToken(ctx context.Context, req *pb.CheckTokenRequest) (*pb.CheckTokenResponse, error) {
	// Проверяем входные данные
	if req.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}

	// Вызываем сервисный метод
	valid, userID, err := s.service.CheckToken(ctx, req.Token)
	if err != nil {
		// Возвращаем соответствующую ошибку gRPC
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	// Возвращаем результат проверки
	return &pb.CheckTokenResponse{
		Valid:  valid,
		UserId: userID,
	}, nil
}
