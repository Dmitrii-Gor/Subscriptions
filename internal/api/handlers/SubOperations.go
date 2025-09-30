package handlers

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"net/http"
	"test/internal/Models"
	logger "test/pkg"
)

func (h *DbHandler) HealthCheck(c *gin.Context) {
	var greeting string
	err := h.DB.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": greeting})
}

func (h *DbHandler) CreateSubscription(c *gin.Context) {
	var subscription Models.Subscription
	err := c.ShouldBindJSON(&subscription)
	if err != nil {
		logger.Error("CreateSubscription: bind JSON failed: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	query := `
		INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	var id string
	err = h.DB.QueryRow(context.Background(),
		query,
		subscription.ServiceName,
		subscription.Price,
		subscription.UserID,
		subscription.StartDate,
		subscription.EndDate,
	).Scan(&id)
	if err != nil {
		logger.Error("Error occured while inserting new Subscription : " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id": id,
	})
}

func (h *DbHandler) DeleteSubscription(c *gin.Context) {
	var subscription Models.Subscription
	err := c.ShouldBindJSON(&subscription)
	if err != nil {
		logger.Error("DeleteSubscription: bind JSON failed: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `DELETE FROM subscriptions 
       WHERE user_id = $1 AND service_name = $2 AND start_date = $3 AND price = $4
       RETURNING id
       `
	var idDeleted string
	err = h.DB.QueryRow(context.Background(), query,
		subscription.UserID,
		subscription.ServiceName,
		subscription.StartDate,
		subscription.Price,
	).Scan(&idDeleted)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Error("DeleteSubscription: No subscription with userID '" + subscription.UserID + "' found")
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		logger.Error("Error occured while deleting Subscription : " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": idDeleted})
}

func (h *DbHandler) GetSubscription(c *gin.Context) {
	id := c.Param("id")

	query := `
		SELECT id, service_name, price, user_id, start_date, end_date
		FROM subscriptions
		WHERE id = $1
		`

	var subscription Models.Subscription
	err := h.DB.QueryRow(context.Background(), query, id).Scan(
		&subscription.ID,
		&subscription.ServiceName,
		&subscription.Price,
		&subscription.UserID,
		&subscription.StartDate,
		&subscription.EndDate,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":           id,
		"service_name": subscription.ServiceName,
		"price":        subscription.Price,
		"user_id":      subscription.UserID,
		"start_date":   subscription.StartDate,
		"end_date":     subscription.EndDate,
	})
}

func (h *DbHandler) UpdateSubscription(c *gin.Context) {
	id := c.Param("id")

	var subscription Models.Subscription
	if err := c.ShouldBindJSON(&subscription); err != nil {
		logger.Error("UpdateSubscription: bind JSON failed: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `
        UPDATE subscriptions
        SET service_name = $1,
            price = $2,
            user_id = $3,
            start_date = $4,
            end_date = $5
        WHERE id = $6
    `

	cmdTag, err := h.DB.Exec(context.Background(),
		query,
		subscription.ServiceName,
		subscription.Price,
		subscription.UserID,
		subscription.StartDate,
		subscription.EndDate,
		id,
	)
	if err != nil {
		logger.Error("UpdateSubscription: query failed: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if cmdTag.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *DbHandler) GetSubscriptionList(c *gin.Context) {
	userID := c.Query("user_id")

	query := `
		SELECT id, service_name, price, user_id, start_date, end_date
		FROM subscriptions
	`
	var args []interface{}
	if userID != "" {
		query += " WHERE user_id = $1"
		args = append(args, userID)
	}

	rows, err := h.DB.Query(context.Background(), query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var subscriptions []Models.Subscription
	for rows.Next() {
		var s Models.Subscription
		if err := rows.Scan(
			&s.ID,
			&s.ServiceName,
			&s.Price,
			&s.UserID,
			&s.StartDate,
			&s.EndDate,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		subscriptions = append(subscriptions, s)
	}

	c.JSON(http.StatusOK, subscriptions)
}
