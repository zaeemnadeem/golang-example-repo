package interfaces

import (
	"context"

	pb "github.com/zaeemnadeem/golang-example-repo/api/proto/v1/content"
	"github.com/zaeemnadeem/golang-example-repo/internal/content/app"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcHandler struct {
	app *app.ContentService
	pb.UnimplementedContentServiceServer
}

func NewGrpcHandler(app *app.ContentService) *GrpcHandler {
	return &GrpcHandler{
		app: app,
	}
}

func (h *GrpcHandler) CreateContent(ctx context.Context, req *pb.CreateContentRequest) (*pb.CreateContentResponse, error) {
	cType := "UNSPECIFIED"
	switch req.Type {
	case pb.ContentType_CONTENT_TYPE_IMAGE:
		cType = "IMAGE"
	case pb.ContentType_CONTENT_TYPE_VIDEO:
		cType = "VIDEO"
	case pb.ContentType_CONTENT_TYPE_HTML:
		cType = "HTML"
	}

	content, err := h.app.CreateContent(ctx, req.Title, req.Url, cType, int(req.DurationSeconds))
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to create content: %v", err)
	}

	return &pb.CreateContentResponse{
		Content: content.ToProto(),
	}, nil
}

func (h *GrpcHandler) AssignContentToScreen(ctx context.Context, req *pb.AssignContentToScreenRequest) (*pb.AssignContentToScreenResponse, error) {
	err := h.app.AssignContent(ctx, req.ContentId, req.ScreenId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to assign content: %v", err)
	}

	return &pb.AssignContentToScreenResponse{
		Success: true,
	}, nil
}

func (h *GrpcHandler) GetScreenContent(ctx context.Context, req *pb.GetScreenContentRequest) (*pb.GetScreenContentResponse, error) {
	contents, err := h.app.GetScreenContent(ctx, req.ScreenId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch content: %v", err)
	}

	var pbContents []*pb.Content
	for _, c := range contents {
		pbContents = append(pbContents, c.ToProto())
	}

	return &pb.GetScreenContentResponse{
		Contents: pbContents,
	}, nil
}
