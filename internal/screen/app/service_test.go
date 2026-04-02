package app

import (
	"context"
	"testing"

	"github.com/zaeemnadeem/golang-example-repo/internal/screen/domain"
)

// Mock repository
type mockScreenRepo struct {
	createFn func(ctx context.Context, screen *domain.Screen) error
	getFn    func(ctx context.Context, id string) (*domain.Screen, error)
	updateFn func(ctx context.Context, id, status string) error
	migrateFn func() error
}
func (m *mockScreenRepo) Create(ctx context.Context, screen *domain.Screen) error { return m.createFn(ctx, screen) }
func (m *mockScreenRepo) GetByID(ctx context.Context, id string) (*domain.Screen, error) { return m.getFn(ctx, id) }
func (m *mockScreenRepo) UpdateStatus(ctx context.Context, id string, status string) error { return m.updateFn(ctx, id, status) }
func (m *mockScreenRepo) Migrate() error { return m.migrateFn() }

func TestCreateScreen(t *testing.T) {
	tests := []struct {
		name        string
		reqName     string
		reqLocation string
		mockCreate  func(ctx context.Context, screen *domain.Screen) error
		expectErr   bool
	}{
		{
			name:        "Success",
			reqName:     "Lobby Screen",
			reqLocation: "Lobby",
			mockCreate:  func(ctx context.Context, screen *domain.Screen) error { return nil },
			expectErr:   false,
		},
		{
			name:        "Missing Name",
			reqName:     "",
			reqLocation: "Lobby",
			mockCreate:  func(ctx context.Context, screen *domain.Screen) error { return nil },
			expectErr:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := &mockScreenRepo{createFn: tc.mockCreate}
			service := NewScreenService(repo)

			screen, err := service.CreateScreen(context.Background(), tc.reqName, tc.reqLocation)

			if (err != nil) != tc.expectErr {
				t.Fatalf("expected err %v, got %v", tc.expectErr, err)
			}
			if !tc.expectErr && screen == nil {
				t.Fatal("expected screen, got nil")
			}
		})
	}
}
