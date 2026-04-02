package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	pbContent "github.com/zaeemnadeem/golang-example-repo/api/proto/v1/content"
	_ "github.com/zaeemnadeem/golang-example-repo/docs" // Ignore error locally if swagger isn't generated yet
	"github.com/zaeemnadeem/golang-example-repo/internal/screen/app"
	"github.com/zaeemnadeem/golang-example-repo/internal/screen/infrastructure"
	"github.com/zaeemnadeem/golang-example-repo/internal/screen/interfaces"
	"github.com/zaeemnadeem/golang-example-repo/pkg/config"
	"github.com/zaeemnadeem/golang-example-repo/pkg/db"
	"github.com/zaeemnadeem/golang-example-repo/pkg/graceful"
	"github.com/zaeemnadeem/golang-example-repo/pkg/grpc"
	"github.com/zaeemnadeem/golang-example-repo/pkg/logger"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

// @title           Signage Microservices API
// @version         1.0.0
// @description     Public HTTP API for the Screen Service.

// @host      localhost:8001
// @BasePath  /
func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	log := logger.Init(cfg.Env)
	defer log.Sync()

	database, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Could not connect to db", zap.Error(err))
	}

	screenRepo := infrastructure.NewPostgresScreenRepository(database)
	if err := screenRepo.Migrate(); err != nil {
		log.Fatal("Could not migrate screen tables", zap.Error(err))
	}

	screenService := app.NewScreenService(screenRepo)

	// Step 1: Connect to content-service via internal gRPC network.
	conn, err := grpc.NewClient(cfg.ContentServiceAddr)
	if err != nil {
		log.Fatal("Failed to connect to content-service", zap.Error(err))
	}
	defer conn.Close()

	contentClient := pbContent.NewContentServiceClient(conn)

	// Step 2: Set up HTTP Handlers.
	httpHandler := interfaces.NewHttpHandler(screenService, contentClient, log)
	mux := http.NewServeMux()
	httpHandler.RegisterRoutes(mux)

	// Register Swagger route
	mux.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://localhost:%d/swagger/doc.json", cfg.ScreenPort)),
	))

	// Step 3: Start HTTP Server.
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.ScreenPort),
		Handler: mux,
	}

	ctx, cancel := graceful.WatchContext()
	defer cancel()

	go func() {
		log.Info("Starting Screen Service HTTP server", zap.Int("port", cfg.ScreenPort))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("HTTP server error", zap.Error(err))
		}
	}()

	<-ctx.Done()
	log.Info("Stopping HTTP server gracefully...")
	server.Shutdown(context.Background())
}
