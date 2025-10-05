package service

import "test/internal/storage"

type SubscriptionService struct {
	storage *storage.SubscriptionStorage
}

func NewSubscriptionService(s *storage.SubscriptionStorage) *SubscriptionService {
	return &SubscriptionService{storage: s}
}
