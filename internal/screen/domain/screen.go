package domain

import (
	"time"

	pb "github.com/zaeemnadeem/golang-example-repo/api/proto/v1/screen"
)

// Screen represents the domain model for a digital signage screen.
type Screen struct {
	ID        string    `gorm:"primaryKey"`
	Name      string    `gorm:"not null"`
	Location  string    `gorm:"not null"`
	Status    string    `gorm:"not null;default:'OFFLINE'"` // e.g., OFFLINE, ONLINE, MAINTENANCE
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ToProto converts a domain model to a gRPC proto message.
func (s *Screen) ToProto() *pb.Screen {
	status := pb.ScreenStatus_SCREEN_STATUS_UNSPECIFIED
	switch s.Status {
	case "OFFLINE":
		status = pb.ScreenStatus_SCREEN_STATUS_OFFLINE
	case "ONLINE":
		status = pb.ScreenStatus_SCREEN_STATUS_ONLINE
	case "MAINTENANCE":
		status = pb.ScreenStatus_SCREEN_STATUS_MAINTENANCE
	}

	return &pb.Screen{
		Id:        s.ID,
		Name:      s.Name,
		Location:  s.Location,
		Status:    status,
		CreatedAt: s.CreatedAt.Unix(),
		UpdatedAt: s.UpdatedAt.Unix(),
	}
}
