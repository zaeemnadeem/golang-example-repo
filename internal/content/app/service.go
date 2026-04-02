package app

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/zaeemnadeem/golang-example-repo/internal/content/domain"
	"github.com/zaeemnadeem/golang-example-repo/internal/content/infrastructure"
)

type ContentService struct {
	repo infrastructure.ContentRepository
}

func NewContentService(repo infrastructure.ContentRepository) *ContentService {
	return &ContentService{
		repo: repo,
	}
}

func (s *ContentService) CreateContent(ctx context.Context, title, url, cType string, duration int) (*domain.Content, error) {
	if title == "" || url == "" || cType == "" {
		return nil, errors.New("title, url, and type are required")
	}

	content := &domain.Content{
		ID:              uuid.New().String(),
		Title:           title,
		URL:             url,
		Type:            cType,
		DurationSeconds: duration,
		CreatedAt:       time.Now(),
	}

	if err := s.repo.CreateContent(ctx, content); err != nil {
		return nil, err
	}

	return content, nil
}

func (s *ContentService) AssignContent(ctx context.Context, contentID, screenID string) error {
	if contentID == "" || screenID == "" {
		return errors.New("content id and screen id are required")
	}

	assignment := &domain.ScreenContent{
		ScreenID:  screenID,
		ContentID: contentID,
		CreatedAt: time.Now(),
	}

	return s.repo.AssignContentToScreen(ctx, assignment)
}

func (s *ContentService) GetScreenContent(ctx context.Context, screenID string) ([]*domain.Content, error) {
	if screenID == "" {
		return nil, errors.New("screen id is required")
	}

	return s.repo.GetContentsForScreen(ctx, screenID)
}
