package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"habit-tracker/internal/middleware"
	"habit-tracker/internal/service"
)

type AchievementHandler struct {
	svc *service.AchievementService
}

func NewAchievementHandler(svc *service.AchievementService) *AchievementHandler {
	return &AchievementHandler{svc: svc}
}

func (h *AchievementHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("", h.ListAll)
	r.GET("/user", h.ListUserAchievements)
}

func (h *AchievementHandler) ListAll(c *gin.Context) {
	list, err := h.svc.ListAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}

func (h *AchievementHandler) ListUserAchievements(c *gin.Context) {
	uid, ok := c.Get(middleware.ContextUserIDKey)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	list, err := h.svc.ListByUser(c.Request.Context(), uid.(uint64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}
