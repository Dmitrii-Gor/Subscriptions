package handlers

import (
	"go.uber.org/zap"
	"net/http"

	"github.com/gin-gonic/gin"
	"test/internal/models"
	logger "test/pkg"
)

func (h *SubscriptionHandler) HealthCheck(c *gin.Context) {
	if err := h.repo.HealthCheck(c.Request.Context()); err != nil {
		logger.Error("HealthCheck: storage unavailable", zap.Error(err))
		sendErrorResponse(c, http.StatusServiceUnavailable, "Storage unavailable")
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *SubscriptionHandler) CreateSubscription(c *gin.Context) {
	var sub models.Subscription
	if err := c.ShouldBindJSON(&sub); err != nil {
		logger.Error("CreateSubscription: bind JSON failed: " + err.Error())
		sendErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	id, err := h.repo.Create(c.Request.Context(), &sub)
	if err != nil {
		logger.Error("CreateSubscription: DB error: " + err.Error())
		sendErrorResponse(c, http.StatusInternalServerError, "failed in create subscription")
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *SubscriptionHandler) GetSubscription(c *gin.Context) {
	id := c.Param("id")

	sub, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		logger.Error("GetSubscription: DB error: " + err.Error())
		sendErrorResponse(c, http.StatusInternalServerError, "failed in get subscription")
		return
	}

	c.JSON(http.StatusOK, sub)
}

func (h *SubscriptionHandler) DeleteSubscription(c *gin.Context) {
	userID := c.Query("user_id")
	serviceName := c.Query("service_name")

	if userID == "" || serviceName == "" {
		logger.Error("DeleteSubscription: invalid request body")
		sendErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.repo.Delete(c.Request.Context(), userID, serviceName); err != nil {
		logger.Error("DeleteSubscription: " + err.Error())
		sendErrorResponse(c, http.StatusInternalServerError, "failed in delete subscription")
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

func (h *SubscriptionHandler) GetSubscriptionList(c *gin.Context) {
	userID := c.Query("user_id")

	list, err := h.repo.List(c.Request.Context(), userID)
	if err != nil {
		logger.Error("GetSubscriptionList: DB error: " + err.Error())
		sendErrorResponse(c, http.StatusInternalServerError, "failed in get subscription list")
		return
	}

	c.JSON(http.StatusOK, list)
}
