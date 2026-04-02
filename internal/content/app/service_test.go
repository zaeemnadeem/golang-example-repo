package app

import (
	"context"
	"testing"

	"github.com/zaeemnadeem/golang-example-repo/internal/content/domain"
)

// Mock repo
type mockContentRepo struct {
	createFn func(ctx context.Context, content *domain.Content) error
	assignFn func(ctx context.Context, assignment *domain.ScreenContent) error
	getFn    func(ctx context.Context, screenID string) ([]*domain.Content, error)
	migrateFn func() error
}
func (m *mockContentRepo) CreateContent(ctx context.Context, c *domain.Content) error { return m.createFn(ctx, c) }
func (m *mockContentRepo) AssignContentToScreen(ctx context.Context, a *domain.ScreenContent) error { return m.assignFn(ctx, a) }
func (m *mockContentRepo) GetContentsForScreen(ctx context.Context, id string) ([]*domain.Content, error) { return m.getFn(ctx, id) }
func (m *mockContentRepo) Migrate() error { return m.migrateFn() }

func TestCreateContent(t *testing.T) {
	tests := []struct {
		name       string
		reqTitle   string
		reqURL     string
		reqType    string
		reqDur     int
		mockCreate func(ctx context.Context, content *domain.Content) error
		expectErr  bool
	}{
		{
			name:       "Success Image",
			reqTitle:   "Promo",
			reqURL:     "http://example.com/promo.jpg",
			reqType:    "IMAGE",
			reqDur:     10,
			mockCreate: func(ctx context.Context, content *domain.Content) error { return nil },
			expectErr:  false,
		},
		{
			name:       "Missing Type",
			reqTitle:   "Promo",
			reqURL:     "http://example.com",
			reqType:    "",
			reqDur:     10,
			mockCreate: func(ctx context.Context, content *domain.Content) error { return nil },
			expectErr:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := &mockContentRepo{createFn: tc.mockCreate}
			service := NewContentService(repo)

			c, err := service.CreateContent(context.Background(), tc.reqTitle, tc.reqURL, tc.reqType, tc.reqDur)

			if (err != nil) != tc.expectErr {
				t.Fatalf("expected err %v, got %v", tc.expectErr, err)
			}
			if !tc.expectErr && c == nil {
				t.Fatal("expected content, got nil")
			}
		})
	}
}
