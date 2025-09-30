package handlers

import "github.com/jackc/pgx/v5/pgxpool"

type DbHandler struct {
	DB *pgxpool.Pool
}

func NewHandler(db *pgxpool.Pool) *DbHandler {
	return &DbHandler{DB: db}
}
