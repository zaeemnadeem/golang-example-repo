package grpc

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnaryLoggingInterceptor logs incoming gRPC requests.
func UnaryLoggingInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		// Call the handler
		resp, err := handler(ctx, req)

		duration := time.Since(start)

		if err != nil {
			st, _ := status.FromError(err)
			logger.Error("gRPC Error",
				zap.String("method", info.FullMethod),
				zap.Duration("duration", duration),
				zap.String("code", st.Code().String()),
				zap.Error(err),
			)
		} else {
			logger.Info("gRPC Success",
				zap.String("method", info.FullMethod),
				zap.Duration("duration", duration),
				zap.String("code", codes.OK.String()),
			)
		}

		return resp, err
	}
}
