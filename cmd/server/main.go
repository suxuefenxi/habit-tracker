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
	checkinRepo := repository.NewCheckinRepository(db.DB)
	pointsRepo := repository.NewPointsRepository(db.DB)
	achRepo := repository.NewAchievementRepository(db.DB)
	userAchRepo := repository.NewUserAchievementRepository(db.DB)

	jwtManager := utils.NewJWTManager(cfg.JWTSecret, 0)
	authSvc := service.NewAuthService(userRepo, jwtManager)
	authHandler := handler.NewAuthHandler(authSvc)

	pointsSvc := service.NewPointsService(userRepo, pointsRepo)
	achSvc := service.NewAchievementService(achRepo, userAchRepo)
	habitSvc := service.NewHabitService(habitRepo)
	checkinSvc := service.NewCheckinService(habitRepo, userRepo, checkinRepo, pointsSvc, achSvc)
	leaderboardSvc := service.NewLeaderboardService(userRepo, pointsRepo)
	userStatsSvc := service.NewUserStatsService(userRepo, habitRepo, checkinRepo, pointsSvc)
	habitHandler := handler.NewHabitHandler(habitSvc)
	checkinHandler := handler.NewCheckinHandler(checkinSvc)
	authMW := middleware.AuthMiddleware(jwtManager)
	leaderboardHandler := handler.NewLeaderboardHandler(leaderboardSvc)
	userHandler := handler.NewUserHandler(userStatsSvc)

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	router.Register(r, router.Deps{
		AuthHandler:        authHandler,
		HabitHandler:       habitHandler,
		CheckinHandler:     checkinHandler,
		LeaderboardHandler: leaderboardHandler,
		UserHandler:        userHandler,
		AuthMW:             authMW,
	})

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
