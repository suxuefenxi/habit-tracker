package router

import (
	"habit-tracker/internal/handler"

	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine, authHandler *handler.AuthHandler) {
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	api := r.Group("/api/v1")
	api.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	authGroup := api.Group("/auth")
	authHandler.RegisterRoutes(authGroup)
}
