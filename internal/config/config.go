// 读取环境变量/配置
package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	DBDSN     string
	Port      string
	JWTSecret string
}

func Load() (Config, error) {
	var cfg Config

	cfg.DBDSN = strings.TrimSpace(os.Getenv("DB_DSN"))
	if cfg.DBDSN == "" {
		return Config{}, fmt.Errorf("missing env DB_DSN")
	}

	cfg.Port = strings.TrimSpace(os.Getenv("PORT"))
	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	cfg.JWTSecret = strings.TrimSpace(os.Getenv("JWT_SECRET"))
	if cfg.JWTSecret == "" {
		return Config{}, fmt.Errorf("missing env JWT_SECRET")
	}

	return cfg, nil
}
