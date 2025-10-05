package models

import "time"

type Subscription struct {
	ID          string     `json:"id,omitempty"`
	ServiceName string     `json:"service_name" binding:"required"`
	Price       int        `json:"price" binding:"required,gte=0"`
	UserID      string     `json:"user_id" binding:"required,uuid"`
	StartDate   time.Time  `json:"start_date" binding:"required"`
	EndDate     *time.Time `json:"end_date,omitempty"`
}
