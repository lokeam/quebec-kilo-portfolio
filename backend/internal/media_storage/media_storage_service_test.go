package media_storage

import (
	"context"
	"errors"
	"testing"

	"github.com/lokeam/qko-beta/internal/analytics"
	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAnalyticsService implements the analytics.Service interface for testing
type MockAnalyticsService struct {
	mock.Mock
}

func (m *MockAnalyticsService) GetAnalytics(ctx context.Context, userID string, domains []string) (map[string]any, error) {
	args := m.Called(ctx, userID, domains)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]any), args.Error(1)
}

func (m *MockAnalyticsService) GetGeneralStats(ctx context.Context, userID string) (*analytics.GeneralStats, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*analytics.GeneralStats), args.Error(1)
}

func (m *MockAnalyticsService) GetFinancialStats(ctx context.Context, userID string) (*analytics.FinancialStats, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*analytics.FinancialStats), args.Error(1)
}

func (m *MockAnalyticsService) GetStorageStats(ctx context.Context, userID string) (*analytics.StorageStats, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*analytics.StorageStats), args.Error(1)
}

func (m *MockAnalyticsService) GetInventoryStats(ctx context.Context, userID string) (*analytics.InventoryStats, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*analytics.InventoryStats), args.Error(1)
}

func (m *MockAnalyticsService) GetWishlistStats(ctx context.Context, userID string) (*analytics.WishlistStats, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*analytics.WishlistStats), args.Error(1)
}

func (m *MockAnalyticsService) InvalidateDomain(ctx context.Context, userID string, domain string) error {
	args := m.Called(ctx, userID, domain)
	return args.Error(0)
}

func (m *MockAnalyticsService) InvalidateDomains(ctx context.Context, userID string, domains []string) error {
	args := m.Called(ctx, userID, domains)
	return args.Error(0)
}

func TestNewMediaStorageService(t *testing.T) {
	ctx := context.Background()
	mockLogger := testutils.NewTestLogger()

	// Create empty storage stats for test connection
	emptyStats := &analytics.StorageStats{}

	tests := []struct {
		name            string
		analyticsService analytics.Service
		logger          interfaces.Logger
		mockSetup       func(*MockAnalyticsService)
		wantErr         bool
		expectedErr     error
	}{
		{
			name:            "successful creation",
			analyticsService: new(MockAnalyticsService),
			logger:          mockLogger,
			mockSetup: func(mockAnalytics *MockAnalyticsService) {
				mockAnalytics.On("GetStorageStats", ctx, "test-connection").
					Return(emptyStats, nil)
			},
			wantErr: false,
		},
		{
			name:            "nil analytics service",
			analyticsService: nil,
			logger:          mockLogger,
			mockSetup:       nil,
			wantErr:         true,
			expectedErr:     ErrInvalidInput,
		},
		{
			name:            "nil logger",
			analyticsService: new(MockAnalyticsService),
			logger:          nil,
			mockSetup:       nil,
			wantErr:         true,
			expectedErr:     ErrInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup(tt.analyticsService.(*MockAnalyticsService))
			}

			service, err := NewMediaStorageService(tt.analyticsService, tt.logger)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, service)
				if tt.expectedErr != nil {
					assert.ErrorIs(t, err, tt.expectedErr)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, service)
			}

			if mockAnalytics, ok := tt.analyticsService.(*MockAnalyticsService); ok {
				mockAnalytics.AssertExpectations(t)
			}
		})
	}
}

func TestGetStorageStats(t *testing.T) {
	ctx := context.Background()
	mockLogger := testutils.NewTestLogger()

	// Create empty storage stats for test connection
	emptyStats := &analytics.StorageStats{}

	// Create test storage stats
	testStats := &analytics.StorageStats{
		TotalPhysicalLocations: 5,
		TotalDigitalLocations:  3,
		PhysicalLocations: []analytics.LocationSummary{
			{
				ID:        "loc1",
				Name:      "Physical Location 1",
				ItemCount: 10,
			},
		},
		DigitalLocations: []analytics.LocationSummary{
			{
				ID:        "dig1",
				Name:      "Digital Location 1",
				ItemCount: 5,
			},
		},
	}

	tests := []struct {
		name            string
		userID          string
		mockSetup       func(*MockAnalyticsService)
		expectedStats   *analytics.StorageStats
		expectedErr     error
	}{
		{
			name:   "successful retrieval",
			userID: "user1",
			mockSetup: func(mockAnalytics *MockAnalyticsService) {
				// Setup test connection call
				mockAnalytics.On("GetStorageStats", ctx, "test-connection").
					Return(emptyStats, nil)
				// Setup actual test call
				mockAnalytics.On("GetStorageStats", ctx, "user1").
					Return(testStats, nil)
			},
			expectedStats: testStats,
			expectedErr:   nil,
		},
		{
			name:   "empty user ID",
			userID: "",
			mockSetup: func(mockAnalytics *MockAnalyticsService) {
				// Setup test connection call
				mockAnalytics.On("GetStorageStats", ctx, "test-connection").
					Return(emptyStats, nil)
			},
			expectedStats: nil,
			expectedErr:   ErrInvalidUserID,
		},
		{
			name:   "analytics service unavailable",
			userID: "user1",
			mockSetup: func(mockAnalytics *MockAnalyticsService) {
				// Setup test connection call
				mockAnalytics.On("GetStorageStats", ctx, "test-connection").
					Return(emptyStats, nil)
				// Setup actual test call
				mockAnalytics.On("GetStorageStats", ctx, "user1").
					Return(nil, errors.New("analytics service unavailable"))
			},
			expectedStats: nil,
			expectedErr:   ErrAnalyticsServiceUnavailable,
		},
		{
			name:   "storage stats not found",
			userID: "user1",
			mockSetup: func(mockAnalytics *MockAnalyticsService) {
				// Setup test connection call
				mockAnalytics.On("GetStorageStats", ctx, "test-connection").
					Return(emptyStats, nil)
				// Setup actual test call
				mockAnalytics.On("GetStorageStats", ctx, "user1").
					Return(nil, errors.New("storage stats not found"))
			},
			expectedStats: nil,
			expectedErr:   ErrStorageStatsNotFound,
		},
		{
			name:   "nil stats returned",
			userID: "user1",
			mockSetup: func(mockAnalytics *MockAnalyticsService) {
				// Setup test connection call
				mockAnalytics.On("GetStorageStats", ctx, "test-connection").
					Return(emptyStats, nil)
				// Setup actual test call
				mockAnalytics.On("GetStorageStats", ctx, "user1").
					Return(nil, nil)
			},
			expectedStats: nil,
			expectedErr:   ErrStorageStatsNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock analytics service
			mockAnalytics := new(MockAnalyticsService)
			if tt.mockSetup != nil {
				tt.mockSetup(mockAnalytics)
			}

			// Create service
			service, err := NewMediaStorageService(mockAnalytics, mockLogger)
			assert.NoError(t, err)

			// Execute test
			stats, err := service.GetStorageStats(ctx, tt.userID)

			// Verify results
			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.expectedErr)
				assert.Nil(t, stats)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStats, stats)
			}

			// Verify mock expectations
			mockAnalytics.AssertExpectations(t)
		})
	}
}