package grpc

import (
	"fmt"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// ServerWrapper encapsulates common gRPC server construction.
type ServerWrapper struct {
	Server *grpc.Server
	port   int
	logger *zap.Logger
}

// NewServer creates a new ServerWrapper with base interceptors.
func NewServer(port int, logger *zap.Logger) *ServerWrapper {
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(UnaryLoggingInterceptor(logger)),
	)

	// Enable reflection for easy debugging (e.g., using grpcurl)
	reflection.Register(grpcServer)

	return &ServerWrapper{
		Server: grpcServer,
		port:   port,
		logger: logger,
	}
}

// Start listens and serves on the specified port.
func (s *ServerWrapper) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		s.logger.Error("failed to listen", zap.Error(err))
		return err
	}

	s.logger.Info("gRPC server is running", zap.Int("port", s.port))
	return s.Server.Serve(lis)
}

// Stop gracefully shuts down the gRPC server.
func (s *ServerWrapper) Stop() {
	s.logger.Info("Stopping gRPC server gracefully...")
	s.Server.GracefulStop()
}
