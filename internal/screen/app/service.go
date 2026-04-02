package app

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/zaeemnadeem/golang-example-repo/internal/screen/domain"
	"github.com/zaeemnadeem/golang-example-repo/internal/screen/infrastructure"
)

// ScreenService defines the business logic for standard screen operations.
type ScreenService struct {
	repo infrastructure.ScreenRepository
}

func NewScreenService(repo infrastructure.ScreenRepository) *ScreenService {
	return &ScreenService{
		repo: repo,
	}
}

func (s *ScreenService) CreateScreen(ctx context.Context, name, location string) (*domain.Screen, error) {
	if name == "" || location == "" {
		return nil, errors.New("name and location must be provided")
	}

	screen := &domain.Screen{
		ID:        uuid.New().String(),
		Name:      name,
		Location:  location,
		Status:    "OFFLINE", // Default
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, screen); err != nil {
		return nil, err
	}

	return screen, nil
}

func (s *ScreenService) GetScreen(ctx context.Context, id string) (*domain.Screen, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}

	screen, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if screen == nil {
		return nil, errors.New("screen not found")
	}

	return screen, nil
}

func (s *ScreenService) UpdateScreenStatus(ctx context.Context, id, status string) (*domain.Screen, error) {
	if id == "" || status == "" {
		return nil, errors.New("id and status required")
	}

	// Simple validation
	if status != "ONLINE" && status != "OFFLINE" && status != "MAINTENANCE" {
		return nil, errors.New("invalid status")
	}

	if err := s.repo.UpdateStatus(ctx, id, status); err != nil {
		return nil, err
	}

	// Always return the updated screen
	return s.GetScreen(ctx, id)
}
