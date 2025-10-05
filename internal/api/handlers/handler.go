package handlers

import (
	"test/internal/storage"
)

type SubscriptionHandler struct {
	repo storage.SubscriptionRepository
}

func NewSubscriptionHandler(repo storage.SubscriptionRepository) *SubscriptionHandler {
	return &SubscriptionHandler{repo: repo}
}
