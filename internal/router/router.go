package router

import (
	"habit-tracker/internal/handler"

	"github.com/gin-gonic/gin"
)

type Deps struct {
	AuthHandler    *handler.AuthHandler
	HabitHandler   *handler.HabitHandler
	CheckinHandler *handler.CheckinHandler
	AuthMW         gin.HandlerFunc
}

func Register(r *gin.Engine, deps Deps) {
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	api := r.Group("/api/v1")
	api.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	authGroup := api.Group("/auth")
	deps.AuthHandler.RegisterRoutes(authGroup)

	habits := api.Group("/habits")
	habits.Use(deps.AuthMW)
	deps.HabitHandler.RegisterRoutes(habits)

	checkins := api.Group("")
	checkins.Use(deps.AuthMW)
	deps.CheckinHandler.RegisterRoutes(checkins)
}
