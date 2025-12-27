package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"habit-tracker/internal/middleware"
	"habit-tracker/internal/service"
	"habit-tracker/internal/utils"
)

type HabitHandler struct {
	habitSvc *service.HabitService
}

func NewHabitHandler(habitSvc *service.HabitService) *HabitHandler {
	return &HabitHandler{habitSvc: habitSvc}
}

type habitResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func writeHabitOK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, habitResponse{Code: 0, Message: "ok", Data: data})
}

func writeHabitError(c *gin.Context, status int, msg string) {
	c.JSON(status, habitResponse{Code: 1, Message: msg})
}

type createHabitRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	TargetType  string `json:"target_type" binding:"required"`
	TargetTimes int    `json:"target_times" binding:"required"`
	StartDate   string `json:"start_date" binding:"required"`
}

type updateHabitRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	TargetType  string `json:"target_type" binding:"required"`
	TargetTimes int    `json:"target_times" binding:"required"`
	StartDate   string `json:"start_date" binding:"required"`
	IsActive    *bool  `json:"is_active"`
}

type toggleHabitStatusRequest struct {
	IsActive bool `json:"is_active" binding:"required"`
}

func (h *HabitHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("", h.ListHabits)
	rg.POST("", h.CreateHabit)
	rg.GET(":id", h.GetHabit)
	rg.PUT(":id", h.UpdateHabit)
	rg.PATCH(":id/status", h.ToggleHabitStatus)
}

func (h *HabitHandler) CreateHabit(c *gin.Context) {
	uid, ok := c.Get(middleware.ContextUserIDKey)
	if !ok {
		writeHabitError(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req createHabitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeHabitError(c, http.StatusBadRequest, "invalid request")
		return
	}
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		writeHabitError(c, http.StatusBadRequest, "invalid start_date")
		return
	}

	habit, err := h.habitSvc.Create(c.Request.Context(), uid.(uint64), service.HabitInput{
		Name:        req.Name,
		Description: req.Description,
		TargetType:  req.TargetType,
		TargetTimes: req.TargetTimes,
		StartDate:   startDate,
	})
	if err != nil {
		writeHabitError(c, http.StatusBadRequest, err.Error())
		return
	}
	writeHabitOK(c, habit)
}

func (h *HabitHandler) UpdateHabit(c *gin.Context) {
	uid, ok := c.Get(middleware.ContextUserIDKey)
	if !ok {
		writeHabitError(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	habitID, err := utils.ParseIDParam(c.Param("id"))
	if err != nil {
		writeHabitError(c, http.StatusBadRequest, "invalid id")
		return
	}

	var req updateHabitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeHabitError(c, http.StatusBadRequest, "invalid request")
		return
	}
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		writeHabitError(c, http.StatusBadRequest, "invalid start_date")
		return
	}

	habit, err := h.habitSvc.Update(c.Request.Context(), uid.(uint64), habitID, service.HabitInput{
		Name:        req.Name,
		Description: req.Description,
		TargetType:  req.TargetType,
		TargetTimes: req.TargetTimes,
		StartDate:   startDate,
		IsActive:    req.IsActive,
	})
	if err != nil {
		writeHabitError(c, statusFromHabitError(err), err.Error())
		return
	}
	writeHabitOK(c, habit)
}

func (h *HabitHandler) ListHabits(c *gin.Context) {
	uid, ok := c.Get(middleware.ContextUserIDKey)
	if !ok {
		writeHabitError(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	var activePtr *bool
	if v := c.Query("is_active"); v != "" {
		parsed, err := strconv.ParseBool(v)
		if err != nil {
			writeHabitError(c, http.StatusBadRequest, "invalid is_active")
			return
		}
		activePtr = &parsed
	}
	habits, err := h.habitSvc.List(c.Request.Context(), uid.(uint64), activePtr)
	if err != nil {
		writeHabitError(c, http.StatusInternalServerError, err.Error())
		return
	}
	writeHabitOK(c, habits)
}

func (h *HabitHandler) GetHabit(c *gin.Context) {
	uid, ok := c.Get(middleware.ContextUserIDKey)
	if !ok {
		writeHabitError(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	habitID, err := utils.ParseIDParam(c.Param("id"))
	if err != nil {
		writeHabitError(c, http.StatusBadRequest, "invalid id")
		return
	}
	habit, err := h.habitSvc.Get(c.Request.Context(), uid.(uint64), habitID)
	if err != nil {
		writeHabitError(c, statusFromHabitError(err), err.Error())
		return
	}
	writeHabitOK(c, habit)
}

func (h *HabitHandler) ToggleHabitStatus(c *gin.Context) {
	uid, ok := c.Get(middleware.ContextUserIDKey)
	if !ok {
		writeHabitError(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	habitID, err := utils.ParseIDParam(c.Param("id"))
	if err != nil {
		writeHabitError(c, http.StatusBadRequest, "invalid id")
		return
	}
	var req toggleHabitStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeHabitError(c, http.StatusBadRequest, "invalid request")
		return
	}
	if err := h.habitSvc.SetActive(c.Request.Context(), uid.(uint64), habitID, req.IsActive); err != nil {
		writeHabitError(c, statusFromHabitError(err), err.Error())
		return
	}
	writeHabitOK(c, gin.H{"is_active": req.IsActive})
}

func statusFromHabitError(err error) int {
	switch {
	case errors.Is(err, service.ErrHabitForbidden):
		return http.StatusForbidden
	case errors.Is(err, service.ErrHabitNotFound):
		return http.StatusNotFound
	default:
		return http.StatusBadRequest
	}
}
