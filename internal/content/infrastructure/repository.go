package infrastructure

import (
	"context"

	"github.com/zaeemnadeem/golang-example-repo/internal/content/domain"
	"gorm.io/gorm"
)

// ContentRepository defines data access methods for content scheduling and storage.
type ContentRepository interface {
	CreateContent(ctx context.Context, content *domain.Content) error
	AssignContentToScreen(ctx context.Context, assignment *domain.ScreenContent) error
	GetContentsForScreen(ctx context.Context, screenID string) ([]*domain.Content, error)
	Migrate() error
}

type postgresContentRepo struct {
	db *gorm.DB
}

// NewPostgresContentRepository creates a new repo constraint object.
func NewPostgresContentRepository(db *gorm.DB) ContentRepository {
	return &postgresContentRepo{db: db}
}

func (r *postgresContentRepo) Migrate() error {
	return r.db.AutoMigrate(&domain.Content{}, &domain.ScreenContent{})
}

func (r *postgresContentRepo) CreateContent(ctx context.Context, content *domain.Content) error {
	return r.db.WithContext(ctx).Create(content).Error
}

func (r *postgresContentRepo) AssignContentToScreen(ctx context.Context, assignment *domain.ScreenContent) error {
	return r.db.WithContext(ctx).Create(assignment).Error
}

func (r *postgresContentRepo) GetContentsForScreen(ctx context.Context, screenID string) ([]*domain.Content, error) {
	var contents []*domain.Content
	
	err := r.db.WithContext(ctx).
		Table("contents").
		Joins("JOIN screen_contents ON screen_contents.content_id = contents.id").
		Where("screen_contents.screen_id = ?", screenID).
		Find(&contents).Error

	return contents, err
}
