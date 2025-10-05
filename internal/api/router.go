package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"test/internal/api/handlers"
	"test/internal/api/middleware"
	"test/internal/storage"
)

func InitGinRouter(pool *pgxpool.Pool) (*gin.Engine, error) {
	r := gin.New()
	r.Use(middleware.GinLogger(), middleware.GinRecovery())
	if err := r.SetTrustedProxies(nil); err != nil {
		return nil, fmt.Errorf("set trusted proxies: %v", err)
	}

	subRepo := storage.NewSubscriptionStorage(pool)
	subHandler := handlers.NewSubscriptionHandler(subRepo)

	api := r.Group("/api/v1")
	api.GET("/health", subHandler.HealthCheck)
	api.POST("/create-sub", subHandler.CreateSubscription)
	api.GET("/get-sub/:id", subHandler.GetSubscription)
	api.DELETE("/delete-sub", subHandler.DeleteSubscription)
	api.GET("/get-sub-list", subHandler.GetSubscriptionList)

	return r, nil
}
