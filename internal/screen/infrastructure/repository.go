package infrastructure

import (
	"context"

	"github.com/zaeemnadeem/golang-example-repo/internal/screen/domain"
	"gorm.io/gorm"
)

// ScreenRepository defines the data access methods for screens.
type ScreenRepository interface {
	Create(ctx context.Context, screen *domain.Screen) error
	GetByID(ctx context.Context, id string) (*domain.Screen, error)
	UpdateStatus(ctx context.Context, id string, status string) error
	Migrate() error
}

type postgresScreenRepo struct {
	db *gorm.DB
}

// NewPostgresScreenRepository instantiates a new postgres repo for screens.
func NewPostgresScreenRepository(db *gorm.DB) ScreenRepository {
	return &postgresScreenRepo{db: db}
}

func (r *postgresScreenRepo) Migrate() error {
	return r.db.AutoMigrate(&domain.Screen{})
}

func (r *postgresScreenRepo) Create(ctx context.Context, screen *domain.Screen) error {
	return r.db.WithContext(ctx).Create(screen).Error
}

func (r *postgresScreenRepo) GetByID(ctx context.Context, id string) (*domain.Screen, error) {
	var screen domain.Screen
	if err := r.db.WithContext(ctx).First(&screen, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Return nil on not found
		}
		return nil, err
	}
	return &screen, nil
}

func (r *postgresScreenRepo) UpdateStatus(ctx context.Context, id string, status string) error {
	return r.db.WithContext(ctx).Model(&domain.Screen{}).Where("id = ?", id).Update("status", status).Error
}
