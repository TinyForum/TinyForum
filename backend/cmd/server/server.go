package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"tiny-forum/internal/infra/config"
	"tiny-forum/internal/startup"
	"tiny-forum/internal/wire"
	"tiny-forum/pkg/logger"
)

// runApp 运行应用主逻辑
func runApp(configDir string, version string) error {
	// ===== 1. 加载静态配置（用于初始化） =====
	staticCfg, err := loadConfig(configDir)
	if err != nil {
		printConfigError(err)
		os.Exit(1)
	}

	// 加载配置成功
	fmt.Println("Config loaded successfully")

	// 初始化日志
	if err = logger.Init(logger.Config(staticCfg.ToLoggerConfig())); err != nil {
		logger.Fatal(fmt.Sprintf("Failed to init logger: %v\n", err))
	}

	// ===== 2. 创建动态配置管理器 =====
	dynCfg, err := config.NewDynamicConfig(configDir)
	logger.Infof("动态配置: %s", dynCfg.Get)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Failed to create dynamic config: %v", err))
		printConfigError(err)
		os.Exit(1)
	}
	defer dynCfg.Close()

	// 打印启动信息
	startup.PrintBanner(version)
	startup.PrintStartupInfo(dynCfg.Get())
	if dynCfg.GetBasic().Log.Level == "debug" {
		startup.PrintConfigSummary(dynCfg.Get())
	}

	// ===== 3. 初始化核心应用（传入动态配置） =====
	app, err := wire.InitAppWithDynamic(dynCfg)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Failed to initialize app: %v\n", err))
	}
	defer logger.CloseDB()

	// ===== 4. 启动 HTTP 服务器 =====
	addr := fmt.Sprintf(":%d", dynCfg.GetBasic().Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: app.Engine,
	}

	go func() {
		logger.Info(fmt.Sprintf("✅ Server is running on http://localhost%s", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal(fmt.Sprintf("Server failed: %v", err))
		}
	}()

	// ===== 5. 等待信号 =====
	return waitForShutdown(srv, app, dynCfg)
}

// waitForShutdown 等待关闭信号并优雅关闭
func waitForShutdown(srv *http.Server, app *wire.App, dynCfg *config.DynamicConfig) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	for {
		sig := <-quit
		switch sig {
		case syscall.SIGHUP:
			// 手动重新加载配置
			logger.Info("Received SIGHUP, reloading config...")
			if err := dynCfg.Reload(); err != nil {
				logger.Error(fmt.Sprintf("Failed to reload config: %v", err))
			} else {
				logger.Info("Config reloaded successfully")
			}
		default:
			// 优雅退出
			return gracefulShutdown(srv, app)
		}
	}
}

// gracefulShutdown 优雅关闭服务器
func gracefulShutdown(srv *http.Server, app *wire.App) error {
	logger.Info("Shutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error(fmt.Sprintf("Server forced to shutdown: %v", err))
	}

	// 停止机器人调度器
	if app.BotSvc != nil {
		app.BotSvc.StopScheduler()
		logger.Info("Bot scheduler stopped")
	}

	logger.Info("Server exited")
	return nil
}
