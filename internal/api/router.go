package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"test/internal/api/handlers"
	"test/internal/api/middleware"
)

func InitGinRouter(pool *pgxpool.Pool) (*gin.Engine, error) {
	r := gin.New()
	r.Use(middleware.GinLogger(), middleware.GinRecovery())
	err := r.SetTrustedProxies(nil)
	if err != nil {
		return nil, fmt.Errorf("set trusted proxies error %v", err)
	}

	h := handlers.NewHandler(pool)

	crudlGroup := r.Group("/api/v1")

	crudlGroup.GET("/health", h.HealthCheck)

	crudlGroup.POST("/create-sub", h.CreateSubscription)
	crudlGroup.DELETE("/delete-sub", h.DeleteSubscription)
	crudlGroup.GET("/get-sub/:id", h.GetSubscription)
	crudlGroup.GET("/get-sub-list", h.GetSubscriptionList)
	crudlGroup.PUT("/update-sub/:id", h.UpdateSubscription)

	return r, nil
}
