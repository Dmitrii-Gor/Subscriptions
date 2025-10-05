package storage

import (
	"context"
	"test/internal/Models"
)

func (s *SubscriptionStorage) HealthCheck(ctx context.Context) error {
	var tmp string
	return s.DB.QueryRow(ctx, "select 'ok'").Scan(&tmp)
}

func (s *SubscriptionStorage) Create(ctx context.Context, sub *Models.Subscription) (string, error) {
	query := `INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
	          VALUES ($1, $2, $3, $4, $5) RETURNING id`
	var id string
	err := s.DB.QueryRow(ctx, query,
		sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate).Scan(&id)
	return id, err
}

func (s *SubscriptionStorage) GetByID(ctx context.Context, id string) (*Models.Subscription, error) {
	query := `SELECT id, service_name, price, user_id, start_date, end_date
	          FROM subscriptions WHERE id = $1`
	var sub Models.Subscription
	err := s.DB.QueryRow(ctx, query, id).Scan(
		&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &sub.EndDate)
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

func (s *SubscriptionStorage) Delete(ctx context.Context, userID, serviceName string) error {
	_, err := s.DB.Exec(ctx, `DELETE FROM subscriptions WHERE user_id=$1 AND service_name=$2`, userID, serviceName)
	return err
}

func (s *SubscriptionStorage) List(ctx context.Context, userID string) ([]Models.Subscription, error) {
	rows, err := s.DB.Query(ctx,
		`SELECT id, service_name, price, user_id, start_date, end_date
		 FROM subscriptions WHERE user_id=$1 ORDER BY id`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []Models.Subscription
	for rows.Next() {
		var s Models.Subscription
		if err := rows.Scan(&s.ID, &s.ServiceName, &s.Price, &s.UserID, &s.StartDate, &s.EndDate); err != nil {
			return nil, err
		}
		subs = append(subs, s)
	}
	return subs, nil
}
