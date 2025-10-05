package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"test/internal/models"
	logger "test/pkg"
)

func (h *SubscriptionHandler) HealthCheck(c *gin.Context) {
	if err := h.repo.HealthCheck(c.Request.Context()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *SubscriptionHandler) CreateSubscription(c *gin.Context) {
	var sub models.Subscription
	if err := c.ShouldBindJSON(&sub); err != nil {
		logger.Error("CreateSubscription: bind JSON failed: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.repo.Create(c.Request.Context(), &sub)
	if err != nil {
		logger.Error("CreateSubscription: DB error: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *SubscriptionHandler) GetSubscription(c *gin.Context) {
	id := c.Param("id")

	sub, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sub)
}

func (h *SubscriptionHandler) DeleteSubscription(c *gin.Context) {
	userID := c.Query("user_id")
	serviceName := c.Query("service_name")

	if userID == "" || serviceName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id and service_name required"})
		return
	}

	if err := h.repo.Delete(c.Request.Context(), userID, serviceName); err != nil {
		logger.Error("DeleteSubscription: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

func (h *SubscriptionHandler) GetSubscriptionList(c *gin.Context) {
	userID := c.Query("user_id")

	list, err := h.repo.List(c.Request.Context(), userID)
	if err != nil {
		logger.Error("GetSubscriptionList: DB error: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, list)
}
