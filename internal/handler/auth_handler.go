package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"habit-tracker/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func writeOK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, response{Code: 0, Message: "ok", Data: data})
}

func writeError(c *gin.Context, status int, msg string) {
	c.JSON(status, response{Code: 1, Message: msg})
}

type registerRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Nickname string `json:"nickname"`
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/register", h.Register)
	rg.POST("/login", h.Login)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid request")
		return
	}

	user, err := h.authService.Register(c.Request.Context(), req.Username, req.Password, req.Nickname)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserExists):
			writeError(c, http.StatusConflict, "username already exists")
		default:
			writeError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	writeOK(c, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"nickname": user.Nickname,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid request")
		return
	}

	token, user, err := h.authService.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCredentials):
			writeError(c, http.StatusUnauthorized, "invalid username or password")
		default:
			writeError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	writeOK(c, gin.H{
		"token":    token,
		"user_id":  user.ID,
		"username": user.Username,
		"nickname": user.Nickname,
	})
}
