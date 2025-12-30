package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"habit-tracker/internal/service"
)

type LeaderboardHandler struct {
	svc *service.LeaderboardService
}

func NewLeaderboardHandler(svc *service.LeaderboardService) *LeaderboardHandler {
	return &LeaderboardHandler{svc: svc}
}

type leaderboardResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func writeLeaderboardOK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, leaderboardResponse{Code: 0, Message: "ok", Data: data})
}

func writeLeaderboardError(c *gin.Context, status int, msg string) {
	c.JSON(status, leaderboardResponse{Code: 1, Message: msg})
}

func (h *LeaderboardHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/weekly", h.Weekly)
	rg.GET("/monthly", h.Monthly)
}

func (h *LeaderboardHandler) Weekly(c *gin.Context) {
	entries, err := h.svc.Weekly(c.Request.Context())
	if err != nil {
		writeLeaderboardError(c, http.StatusInternalServerError, "unable to build weekly leaderboard")
		return
	}
	writeLeaderboardOK(c, entries)
}

func (h *LeaderboardHandler) Monthly(c *gin.Context) {
	entries, err := h.svc.Monthly(c.Request.Context())
	if err != nil {
		writeLeaderboardError(c, http.StatusInternalServerError, "unable to build monthly leaderboard")
		return
	}
	writeLeaderboardOK(c, entries)
}
