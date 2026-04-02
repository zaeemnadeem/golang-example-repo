package main

import (
	pb "github.com/zaeemnadeem/golang-example-repo/api/proto/v1/content"
	"github.com/zaeemnadeem/golang-example-repo/internal/content/app"
	"github.com/zaeemnadeem/golang-example-repo/internal/content/infrastructure"
	"github.com/zaeemnadeem/golang-example-repo/internal/content/interfaces"
	"github.com/zaeemnadeem/golang-example-repo/pkg/config"
	"github.com/zaeemnadeem/golang-example-repo/pkg/db"
	"github.com/zaeemnadeem/golang-example-repo/pkg/graceful"
	"github.com/zaeemnadeem/golang-example-repo/pkg/grpc"
	"github.com/zaeemnadeem/golang-example-repo/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

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

	contentRepo := infrastructure.NewPostgresContentRepository(database)
	if err := contentRepo.Migrate(); err != nil {
		log.Fatal("Could not migrate content tables", zap.Error(err))
	}

	contentService := app.NewContentService(contentRepo)
	grpcHandler := interfaces.NewGrpcHandler(contentService)

	server := grpc.NewServer(cfg.ContentPort, log)
	pb.RegisterContentServiceServer(server.Server, grpcHandler)

	// Register gRPC Health Server
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(server.Server, healthServer)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	ctx, cancel := graceful.WatchContext()
	defer cancel()

	go func() {
		if err := server.Start(); err != nil {
			log.Fatal("Failed to start content server", zap.Error(err))
		}
	}()

	<-ctx.Done()
	server.Stop()
}
