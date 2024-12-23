package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/farisarmap/dot-backend-freelance/config"
	"github.com/farisarmap/dot-backend-freelance/internal/adapter"
	"github.com/farisarmap/dot-backend-freelance/internal/handler"
	"github.com/farisarmap/dot-backend-freelance/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	configPath := flag.String("config", "config.json", "Path to config file")
	flag.Parse()

	// Load config
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Init DB
	db, err := config.InitDB(cfg.Database)
	if err != nil {
		log.Fatalf("Error connecting DB: %v", err)
	}

	// Init Redis
	redisClient := config.InitRedis(cfg.Redis)

	userRepo := adapter.NewUserRepository(db)
	orderRepo := adapter.NewOrderRepository(db)

	cacheManager := adapter.NewRedisCache(redisClient, time.Duration(cfg.Redis.TTL)*time.Second)

	userService := service.NewUserService(userRepo, orderRepo, cacheManager)
	orderService := service.NewOrderService(orderRepo, userRepo, cacheManager)

	userHandler := handler.NewUserHandler(userService)
	orderHandler := handler.NewOrderHandler(orderService, cacheManager)

	e := echo.New()

	e.Use(middleware.Recover())

	e.GET("/users", userHandler.GetAllUsers)
	e.POST("/users", userHandler.CreateUser)
	e.GET("/users/:id", userHandler.GetUserByID)
	e.PUT("/users/:id", userHandler.UpdateUser)
	e.PATCH("/users/:id", userHandler.PartialUpdateUser)
	e.DELETE("/users/:id", userHandler.DeleteUser)

	e.POST("/users-and-orders", orderHandler.CreateUserAndOrder)
	e.GET("/orders", orderHandler.GetAllOrders)
	e.POST("/orders", orderHandler.CreateOrder)
	e.GET("/orders/:id", orderHandler.GetOrderByID)
	e.PUT("/orders/:id", orderHandler.UpdateOrder)
	e.PATCH("/orders/:id", orderHandler.PartialUpdateOrder)
	e.DELETE("/orders/:id", orderHandler.DeleteOrder)

	go func() {
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			log.Fatalf("error starting server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
}
