package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"habit-tracker/internal/middleware"
	"habit-tracker/internal/service"
	"habit-tracker/internal/utils"
)

type CheckinHandler struct {
	checkinSvc *service.CheckinService
}

func NewCheckinHandler(checkinSvc *service.CheckinService) *CheckinHandler {
	return &CheckinHandler{checkinSvc: checkinSvc}
}

type checkinResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func writeCheckinOK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, checkinResponse{Code: 0, Message: "ok", Data: data})
}

func writeCheckinError(c *gin.Context, status int, msg string) {
	c.JSON(status, checkinResponse{Code: 1, Message: msg})
}

type checkinRequest struct {
	HabitID  uint64 `json:"habit_id" binding:"required"`
	CountInc int    `json:"count_inc"`
}

type historyQuery struct {
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
}

func (h *CheckinHandler) RegisterRoutes(api *gin.RouterGroup) {
	api.POST("/checkins", h.Checkin)
	api.GET("/habits/:id/checkins", h.ListCheckins)
}

func (h *CheckinHandler) Checkin(c *gin.Context) {
	uidVal, ok := c.Get(middleware.ContextUserIDKey)
	if !ok {
		writeCheckinError(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	userID := uidVal.(uint64)
	var req checkinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeCheckinError(c, http.StatusBadRequest, "invalid request")
		return
	}
	if req.CountInc == 0 {
		req.CountInc = 1
	}
	res, err := h.checkinSvc.Checkin(c.Request.Context(), userID, req.HabitID, req.CountInc)
	if err != nil {
		writeCheckinError(c, statusFromCheckinError(err), err.Error())
		return
	}
	writeCheckinOK(c, res)
}

func (h *CheckinHandler) ListCheckins(c *gin.Context) {
	uidVal, ok := c.Get(middleware.ContextUserIDKey)
	if !ok {
		writeCheckinError(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	userID := uidVal.(uint64)

	habitID, err := utils.ParseIDParam(c.Param("id"))
	if err != nil {
		writeCheckinError(c, http.StatusBadRequest, "invalid id")
		return
	}

	var q historyQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		writeCheckinError(c, http.StatusBadRequest, "invalid query")
		return
	}
	end := time.Now().In(time.Local).Truncate(24 * time.Hour)
	start := end.AddDate(0, 0, -30)
	if q.StartDate != "" {
		if parsed, err := time.Parse("2006-01-02", q.StartDate); err == nil {
			start = parsed
		} else {
			writeCheckinError(c, http.StatusBadRequest, "invalid start_date")
			return
		}
	}
	if q.EndDate != "" {
		if parsed, err := time.Parse("2006-01-02", q.EndDate); err == nil {
			end = parsed
		} else {
			writeCheckinError(c, http.StatusBadRequest, "invalid end_date")
			return
		}
	}

	records, err := h.checkinSvc.ListHistory(c.Request.Context(), userID, habitID, start, end)
	if err != nil {
		writeCheckinError(c, statusFromCheckinError(err), err.Error())
		return
	}
	writeCheckinOK(c, records)
}

func statusFromCheckinError(err error) int {
	switch {
	case errors.Is(err, service.ErrCheckinForbidden):
		return http.StatusForbidden
	case errors.Is(err, service.ErrHabitMissing):
		return http.StatusNotFound
	default:
		return http.StatusBadRequest
	}
}
