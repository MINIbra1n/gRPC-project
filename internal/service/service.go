package service

import (
	"context"
	"errors"
	repo "gRPC-Project/internal/repo"
	"gRPC-Project/pkg/auth"
	"gRPC-Project/pkg/hash"

	"go.uber.org/zap"
)

type AuthService interface {
	Register(ctx context.Context, username, password string) (string, error)
	Login(ctx context.Context, username, password string) (string, string, error)
	CheckToken(ctx context.Context, token string) (bool, string, error)
}

type authService struct {
	repo         repo.RepoUser
	log          *zap.SugaredLogger
	jwtSecretKey string
}

// NewService - конструктор сервиса
func NewService(repo repo.RepoUser, logger *zap.SugaredLogger, jwtSecretKey string) *authService {

	return &authService{
		repo:         repo,
		log:          logger,
		jwtSecretKey: jwtSecretKey,
	}
}

// CreateTask - обработчик запроса на создание задачи
func (s *authService) Register(ctx context.Context, username, password string) (string, error) {
	// Проверяем, существует ли пользователь

	// fmt.Println(username, password)
	_, err := s.repo.GetUserByName(ctx, username)
	if err == nil {
		return "", errors.New("user already exists")
	}

	passwordHash, err := hash.HashPassword(password)
	if err != nil {
		return "", err
	}

	// Создаем пользователя в БД
	user := &repo.User{
		Username: username,

		Password: passwordHash,
	}

	userID, err := s.repo.CreateUser(ctx, *user)
	if err != nil {
		return "", err
	}
	// fmt.Println(userID)
	return userID, nil
}
func (s *authService) Login(ctx context.Context, username, password string) (string, string, error) {
	// Получаем пользователя из БД

	user, err := s.repo.GetUserByName(ctx, username)
	if err != nil {
		return "", "", errors.New("invalid username or password")
	}
	// Проверяем пароль
	if user.Password != password {
		return "", "", errors.New("invalid username or password")
	}

	// Генерируем JWT токен
	token, err := auth.GenerateAccessToken(user.ID, s.jwtSecretKey)
	if err != nil {
		return "", "", err
	}
	// Добавить token в бд!!!
	return token, user.ID, nil
}

func (s *authService) CheckToken(ctx context.Context, token string) (bool, string, error) {
	// Валидируем токен
	claims, err := auth.ValidateAccessToken(token, s.jwtSecretKey)
	if err != nil {
		return false, "", err
	}

	// Проверяем существование пользователя
	_, err = s.repo.GetUserByID(ctx, claims.UserID)
	if err != nil {
		return false, "", errors.New("user not found")
	}

	return true, claims.UserID, nil
}
