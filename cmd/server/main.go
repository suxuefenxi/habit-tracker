// 程序入口
package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"habit-tracker/internal/config"
	"habit-tracker/internal/db"
	"habit-tracker/internal/handler"
	"habit-tracker/internal/middleware"
	"habit-tracker/internal/repository"
	"habit-tracker/internal/router"
	"habit-tracker/internal/service"
	"habit-tracker/internal/utils"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	if _, err := db.Init(cfg.DBDSN); err != nil {
		log.Fatalf("init db: %v", err)
	}

	userRepo := repository.NewUserRepository(db.DB)
	habitRepo := repository.NewHabitRepository(db.DB)
	jwtManager := utils.NewJWTManager(cfg.JWTSecret, 0)
	authSvc := service.NewAuthService(userRepo, jwtManager)
	authHandler := handler.NewAuthHandler(authSvc)
	habitSvc := service.NewHabitService(habitRepo)
	habitHandler := handler.NewHabitHandler(habitSvc)
	authMW := middleware.AuthMiddleware(jwtManager)

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	router.Register(r, authHandler, habitHandler, authMW)

	// 静态文件前端（web/ 目录）
	// - /index.html 等文件可直接访问
	// - 如果你做 SPA，可改为 NoRoute 返回 index.html
	r.StaticFS("/static", http.Dir("./web"))

	addr := ":" + cfg.Port
	log.Printf("listening on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("run server: %v", err)
	}
}
