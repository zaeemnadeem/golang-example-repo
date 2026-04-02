package domain

import (
	"time"

	pb "github.com/zaeemnadeem/golang-example-repo/api/proto/v1/content"
)

// Content represents the domain model for media content.
type Content struct {
	ID              string `gorm:"primaryKey"`
	Title           string `gorm:"not null"`
	URL             string `gorm:"not null"`
	Type            string `gorm:"not null"` // e.g., IMAGE, VIDEO, HTML
	DurationSeconds int
	CreatedAt       time.Time
}

// ScreenContent links a content piece to a screen for scheduling.
type ScreenContent struct {
	ScreenID  string `gorm:"primaryKey"`
	ContentID string `gorm:"primaryKey"`
	CreatedAt time.Time
}

// ToProto converts a domain model to a gRPC proto message.
func (c *Content) ToProto() *pb.Content {
	ctype := pb.ContentType_CONTENT_TYPE_UNSPECIFIED
	switch c.Type {
	case "IMAGE":
		ctype = pb.ContentType_CONTENT_TYPE_IMAGE
	case "VIDEO":
		ctype = pb.ContentType_CONTENT_TYPE_VIDEO
	case "HTML":
		ctype = pb.ContentType_CONTENT_TYPE_HTML
	}

	return &pb.Content{
		Id:              c.ID,
		Title:           c.Title,
		Url:             c.URL,
		Type:            ctype,
		DurationSeconds: int32(c.DurationSeconds),
		CreatedAt:       c.CreatedAt.Unix(),
	}
}
