package grpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// NewClient establishes an insecure gRPC connection to the target address.
func NewClient(target string) (*grpc.ClientConn, error) {
	// In production, use properly secured credentials (TLS).
	return grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
}
