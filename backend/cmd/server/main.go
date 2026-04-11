package main

import (
	"fmt"
	"log"

	"bbs-forum/config"
	"bbs-forum/internal/wire"
	"bbs-forum/pkg/logger"
)

// @title           BBS Forum API
// @version         1.0
// @description     一个基于 Gin 的论坛系统 API
// @host            localhost:8080
// @BasePath        /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and the JWT token.

func main() {
	// Load config
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Init logger
	if err := logger.Init(cfg.Log.Level, cfg.Log.Filename); err != nil {
		log.Fatalf("Failed to init logger: %v", err)
	}

	// Init app (DB + routes)
	app, err := wire.InitApp(cfg)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Failed to initialize app: %v", err))
	}

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	logger.Info(fmt.Sprintf("BBS Forum server starting on %s", addr))

	if err := app.Engine.Run(addr); err != nil {
		logger.Fatal(fmt.Sprintf("Server failed: %v", err))
	}
}
