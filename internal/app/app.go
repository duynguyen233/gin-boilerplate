// Package app configures and runs application.
package app

import (
	"io"
	"sync"
	"template/config"
	controller "template/internal/controller/http"
	"template/internal/repositories"
	"template/internal/services"
	"template/pkg/grpcserver"
	"template/pkg/httpserver"
	"template/pkg/logger"
	"time"

	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	env := os.Getenv("ENV")
	if env == "PROD" {
		err := os.MkdirAll("./.log", os.ModePerm)
		if err != nil {
			panic(err)
		}

		file, err := os.OpenFile(
			"./.log/server.log",
			os.O_APPEND|os.O_CREATE|os.O_WRONLY,
			0664,
		)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		gin.DefaultWriter = io.MultiWriter(os.Stdout, file)

	}
	l := logger.New(cfg.Log.Level)

	//connect postgres with gorm
	db, err := gorm.Open(postgres.Open(cfg.PG.URL), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// Repositories
	userRepo := repositories.NewUserRepository(db)
	// middleware

	// Services
	userService := services.NewUserService(*userRepo)
	// HTTP Server
	handler := gin.New()
	handler.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	controller.NewRouter(handler, l, *userService)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))
	// gRPC Server
	// grpcAdapter := grpcHandler.NewAdapter(couponController, l)
	grpcServer := grpcserver.New(cfg.GRPC.Port)
	// grpcServer.RegisterService(grpcAdapter.RegisterService)
	grpcServer.Start()

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err := <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	case err := <-grpcServer.Notify():
		l.Error(fmt.Errorf("app - Run - grpcServer.Notify: %w", err))
	}

	// Shutdown
	l.Info("Shutting down server...")
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		l.Info("Shutting down gRPC server...")
		if err := grpcServer.Shutdown(); err != nil {
			l.Error(fmt.Errorf("app - Run - grpcServer.Shutdown: %w", err))
		} else {
			l.Info("gRPC server shutdown complete")
		}
	}()
	go func() {
		defer wg.Done()
		l.Info("Shutting down HTTP server...")
		if err := httpServer.Shutdown(); err != nil {
			l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
		} else {
			l.Info("HTTP server shutdown complete")
		}
	}()
	wg.Wait()
	l.Info("All servers have been shut down gracefully")
}
