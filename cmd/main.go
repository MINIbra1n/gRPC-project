package main

import (
	"context"
	"fmt"
	"gRPC-Project/internal/config"
	grpcHandler "gRPC-Project/internal/grpc"
	pb "gRPC-Project/internal/grpc/gen"
	customLogger "gRPC-Project/internal/logger"
	"gRPC-Project/internal/repo"
	"gRPC-Project/internal/service"
	"log"
	"net"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Зафиксил подгрузку .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to load env"))
	}
	// Загружаем конфигурацию из переменных окружения
	var cfg config.AppConfig
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal(errors.Wrap(err, "failed to load configuration"))
	}

	// Инициализация логгера
	logger, err := customLogger.NewLogger(cfg.LogLevel)
	if err != nil {
		log.Fatal(errors.Wrap(err, "error initializing logger"))
	}

	// Подключение к PostgreSQL
	repository, err := repo.NewRepository(context.Background(), cfg.PostgreSQL)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to initialize repository"))
	}

	// Создание сервиса с бизнес-логикой
	serviceInstance := service.NewService(repository, logger, cfg.JWTConfig.SecretKey)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Server.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	// Регистрация обработчиков
	authServer := grpcHandler.NewAuthServer(serviceInstance)
	pb.RegisterAuthServiceServer(grpcServer, authServer)

	// Включаем reflection для удобства отладки (опционально)
	reflection.Register(grpcServer)
	// Добавить Shutdown
	// Запуск сервера
	log.Printf("Starting gRPC server on port %s", cfg.Server.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

// Достать данные из ьтокена
// Добавить метод обновления токенов Acsees and Refresh
