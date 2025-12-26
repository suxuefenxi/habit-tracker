// 程序入口
package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"habit-tracker/internal/config"
	"habit-tracker/internal/db"
	"habit-tracker/internal/router"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	if _, err := db.Init(cfg.DBDSN); err != nil {
		log.Fatalf("init db: %v", err)
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	// 先注册 API 路由（更具体的路由会优先匹配）
	router.Register(r)

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
