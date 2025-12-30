package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"habit-tracker/internal/middleware"
	"habit-tracker/internal/service"
)

type UserHandler struct {
	stats *service.UserStatsService
}

func NewUserHandler(stats *service.UserStatsService) *UserHandler {
	return &UserHandler{stats: stats}
}

type userResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func writeUserOK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, userResponse{Code: 0, Message: "ok", Data: data})
}

func writeUserError(c *gin.Context, status int, msg string) {
	c.JSON(status, userResponse{Code: 1, Message: msg})
}

func (h *UserHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/stats", h.Stats)
}

func (h *UserHandler) Stats(c *gin.Context) {
	uid, ok := c.Get(middleware.ContextUserIDKey)
	if !ok {
		writeUserError(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	stats, err := h.stats.GetStats(c.Request.Context(), uid.(uint64))
	if err != nil {
		writeUserError(c, http.StatusInternalServerError, "unable to load stats")
		return
	}
	writeUserOK(c, stats)
}
