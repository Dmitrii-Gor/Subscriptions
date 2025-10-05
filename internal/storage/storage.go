package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"test/internal/models"
)

type SubscriptionRepository interface {
	HealthCheck(ctx context.Context) error
	Create(ctx context.Context, sub *models.Subscription) (string, error)
	GetByID(ctx context.Context, id string) (*models.Subscription, error)
	Delete(ctx context.Context, userID, serviceName string) error
	List(ctx context.Context, userID string) ([]models.Subscription, error)
}

type SubscriptionStorage struct {
	DB *pgxpool.Pool
}

func NewSubscriptionStorage(db *pgxpool.Pool) *SubscriptionStorage {
	return &SubscriptionStorage{DB: db}
}

func InitStorage() (*pgxpool.Pool, error) {
	dbPool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	err = dbPool.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return dbPool, nil
}
