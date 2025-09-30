package main

import (
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"os"
	"test/internal/api"
	"test/internal/migrations"
	"test/internal/storage"
	"test/pkg"
)

func main() {
	_ = godotenv.Load()

	logger.InitLogger()
	defer logger.Sync()

	pool, err := storage.InitStorage()
	if err != nil {
		logger.Error(err.Error(), zap.Error(err))
		return
	}
	defer pool.Close()

	if os.Getenv("ENV") == "dev" {
		err = migrations.RunMigrations(pool) // Только dev среда
		if err != nil {
			logger.Error(err.Error(), zap.Error(err))
			return
		}
	}

	r, err := api.InitGinRouter(pool)
	if err != nil {
		logger.Error(err.Error(), zap.Error(err))
		return
	}

	if err := r.Run(":8080"); err != nil { // запускаем HTTP сервер
		logger.Error("server run failed", zap.Error(err))
	}
}
