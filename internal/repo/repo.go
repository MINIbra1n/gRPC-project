package repo

import (
	"context"
	"fmt"
	"gRPC-Project/internal/config"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type repoUser struct {
	pool *pgxpool.Pool
}

type RepoUser interface {
	CreateUser(ctx context.Context, user User) (string, error)
	GetUserByName(context.Context, string) (*User, error)
	GetUserByID(ctx context.Context, id string) (*User, error)
}

func NewRepository(ctx context.Context, cfg config.PostgreSQL) (*repoUser, error) {
	// Формируем строку подключения
	connString := fmt.Sprintf(
		`user=%s password=%s host=%s port=%d dbname=%s sslmode=%s 
        pool_max_conns=%d pool_max_conn_lifetime=%s pool_max_conn_idle_time=%s`,
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.SSLMode,
		cfg.PoolMaxConns,
		cfg.PoolMaxConnLifetime.String(),
		cfg.PoolMaxConnIdleTime.String(),
	)

	// Парсим конфигурацию подключения
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse PostgreSQL config")
	}

	// Оптимизация выполнения запросов (кеширование запросов)
	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeCacheDescribe

	// Создаём пул соединений с базой данных
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create PostgreSQL connection pool")
	}

	return &repoUser{pool}, nil
}

func (r *repoUser) CreateUser(ctx context.Context, user User) (string, error) {
	var id string
	now := time.Now()
	user.Created_at = now
	user.Updated_at = now

	query := `
		INSERT INTO users ( username, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	row := r.pool.QueryRow(ctx, query, user.Username, user.Password, user.Created_at, user.Updated_at)
	row.Scan(&id)

	return id, nil
}
func (r *repoUser) GetUserByName(ctx context.Context, userName string) (*User, error) {
	var user User
	query := `SELECT id, username, password, created_at, updated_at FROM users WHERE username = $1`
	row := r.pool.QueryRow(ctx, query, userName)
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Created_at, &user.Updated_at)
	if err != nil {

		return nil, err
	}

	return &user, nil
}
func (r *repoUser) GetUserByID(ctx context.Context, id string) (*User, error) {
	var user User
	query := `SELECT * FROM users WHERE id = $1`
	row := r.pool.QueryRow(ctx, query, id)
	row.Scan(&user)

	return &user, nil
}
