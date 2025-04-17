package digital

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/testutils"
	"github.com/lokeam/qko-beta/internal/testutils/mocks"
)

/*
Behavior:
- GetDigitalLocations retrieves digital locations for a user
  - Attempts to retrieve from cache first
  - Falls back to database on cache miss
  - Caches database results for future requests

- AddDigitalLocation adds a new digital location
  - Validates the location data
  - Adds the location to the database
  - Invalidates the user's cache

- UpdateDigitalLocation updates an existing digital location
  - Validates the location data
  - Updates the location in the database
  - Invalidates both the user's cache and the specific location cache

- DeleteDigitalLocation deletes a digital location
  - Deletes the location from the database
  - Invalidates both the user's cache and the specific location cache

Scenarios:
- GetDigitalLocations:
  - Cache hit
  - Cache miss
  - Database error
  - No locations found

- AddDigitalLocation:
  - Validation success + db success
  - Validation failure
  - Validation success + db failure
  - Cache invalidation failure (should not block success)

- UpdateDigitalLocation:
  - Validation success + db success
  - Validation failure
  - Validation success + db failure
  - Cache invalidation failure (should not block success)

- DeleteDigitalLocation:
  - DB success
  - DB failure
  - Cache invalidation failure (should not block success)
*/

type MockDigitalDbAdapter struct {
	GetUserDigitalLocationsFunc func(ctx context.Context, userID string) ([]models.DigitalLocation, error)
	GetDigitalLocationFunc func(ctx context.Context, userID, locationID string) (models.DigitalLocation, error)
	FindDigitalLocationByNameFunc func(ctx context.Context, userID string, name string) (models.DigitalLocation, error)
	AddDigitalLocationFunc func(ctx context.Context, userID string, location models.DigitalLocation) (models.DigitalLocation, error)
	UpdateDigitalLocationFunc func(ctx context.Context, userID string, location models.DigitalLocation) error
	RemoveDigitalLocationFunc func(ctx context.Context, userID, locationID string) error
}

func (m *MockDigitalDbAdapter) GetUserDigitalLocations(ctx context.Context, userID string) ([]models.DigitalLocation, error) {
	return m.GetUserDigitalLocationsFunc(ctx, userID)
}

func (m *MockDigitalDbAdapter) GetDigitalLocation(ctx context.Context, userID, locationID string) (models.DigitalLocation, error) {
	return m.GetDigitalLocationFunc(ctx, userID, locationID)
}

func (m *MockDigitalDbAdapter) FindDigitalLocationByName(ctx context.Context, userID string, name string) (models.DigitalLocation, error) {
	return m.FindDigitalLocationByNameFunc(ctx, userID, name)
}

func (m *MockDigitalDbAdapter) AddDigitalLocation(ctx context.Context, userID string, location models.DigitalLocation) (models.DigitalLocation, error) {
	return m.AddDigitalLocationFunc(ctx, userID, location)
}

func (m *MockDigitalDbAdapter) UpdateDigitalLocation(ctx context.Context, userID string, location models.DigitalLocation) error {
	return m.UpdateDigitalLocationFunc(ctx, userID, location)
}

func (m *MockDigitalDbAdapter) RemoveDigitalLocation(ctx context.Context, userID, locationID string) error {
	return m.RemoveDigitalLocationFunc(ctx, userID, locationID)
}

func newMockGameDigitalServiceWithDefaults(logger *testutils.TestLogger) *GameDigitalService {
	mockConfig := mocks.NewMockConfig()

	return &GameDigitalService{
		dbAdapter:    mocks.DefaultDigitalDbAdapter(),
		config:       mockConfig,
		logger:       logger,
		validator:    mocks.DefaultDigitalValidator(),
		sanitizer:    mocks.DefaultSanitizer(),
		cacheWrapper: mocks.DefaultDigitalCacheWrapper(),
	}
}

func TestGameDigitalService(t *testing.T) {
	// Set up test logger
	logger := testutils.NewTestLogger()

	// Test cases
	t.Run("GetUserDigitalLocations - Success", func(t *testing.T) {
		// Setup
		service := newMockGameDigitalServiceWithDefaults(logger)

		// Execute
		locations, err := service.GetUserDigitalLocations(context.Background(), "test-user")

		// Verify
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(locations) != 2 {
			t.Errorf("Expected 2 locations, got %d", len(locations))
		}
	})

	t.Run("GetUserDigitalLocations - Error", func(t *testing.T) {
		// Setup
		service := newMockGameDigitalServiceWithDefaults(logger)
		expectedErr := errors.New("test error")

		service.dbAdapter = &MockDigitalDbAdapter{
			GetUserDigitalLocationsFunc: func(ctx context.Context, userID string) ([]models.DigitalLocation, error) {
				return nil, expectedErr
			},
		}

		// Execute
		locations, err := service.GetUserDigitalLocations(context.Background(), "test-user")

		// Verify
		if err == nil {
			t.Error("Expected error, got nil")
		}
		if !errors.Is(err, expectedErr) {
			t.Errorf("Expected error %v, got %v", expectedErr, err)
		}
		if len(locations) != 0 {
			t.Errorf("Expected empty locations slice, got %v", locations)
		}
	})

	t.Run("GetDigitalLocation - Success", func(t *testing.T) {
		// Setup
		service := newMockGameDigitalServiceWithDefaults(logger)

		// Execute
		location, err := service.GetDigitalLocation(context.Background(), "test-user", "test-location")

		// Verify
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if location.ID == "" {
			t.Error("Expected location to be returned")
		}
	})

	t.Run("GetDigitalLocation - Not Found", func(t *testing.T) {
		// Setup
		service := newMockGameDigitalServiceWithDefaults(logger)
		expectedErr := sql.ErrNoRows

		service.dbAdapter = &MockDigitalDbAdapter{
			GetDigitalLocationFunc: func(ctx context.Context, userID, locationID string) (models.DigitalLocation, error) {
				return models.DigitalLocation{}, expectedErr
			},
		}

		// Execute
		location, err := service.GetDigitalLocation(context.Background(), "test-user", "test-location")

		// Verify
		if err == nil {
			t.Error("Expected error, got nil")
		}
		if !errors.Is(err, expectedErr) {
			t.Errorf("Expected error %v, got %v", expectedErr, err)
		}
		if location.ID != "" || location.Name != "" || location.IsActive != false {
			t.Errorf("Expected empty location, got %v", location)
		}
	})

	t.Run("AddDigitalLocation - Success", func(t *testing.T) {
		// Setup
		service := newMockGameDigitalServiceWithDefaults(logger)
		newLocation := models.DigitalLocation{
			Name:     "Test Location",
			IsActive: true,
		}

		// Execute
		createdLocation, err := service.AddDigitalLocation(context.Background(), "test-user", newLocation)

		// Verify
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if createdLocation.Name != newLocation.Name {
			t.Errorf("Expected location name %s, got %s", newLocation.Name, createdLocation.Name)
		}
	})

	t.Run("AddDigitalLocation - Error", func(t *testing.T) {
		// Setup
		service := newMockGameDigitalServiceWithDefaults(logger)
		expectedErr := errors.New("test error")
		newLocation := models.DigitalLocation{
			Name:     "Test Location",
			IsActive: true,
		}

		service.dbAdapter = &MockDigitalDbAdapter{
			AddDigitalLocationFunc: func(ctx context.Context, userID string, location models.DigitalLocation) (models.DigitalLocation, error) {
				return models.DigitalLocation{}, expectedErr
			},
		}

		// Execute
		createdLocation, err := service.AddDigitalLocation(context.Background(), "test-user", newLocation)

		// Verify
		if err == nil {
			t.Error("Expected error, got nil")
		}
		if !errors.Is(err, expectedErr) {
			t.Errorf("Expected error %v, got %v", expectedErr, err)
		}
		if createdLocation.ID != "" || createdLocation.Name != "" || createdLocation.IsActive != false {
			t.Errorf("Expected empty location, got %v", createdLocation)
		}
	})

	t.Run("UpdateDigitalLocation - Success", func(t *testing.T) {
		// Setup
		service := newMockGameDigitalServiceWithDefaults(logger)
		location := models.DigitalLocation{
			ID:       "test-location",
			Name:     "Updated Location",
			IsActive: true,
		}

		// Execute
		err := service.UpdateDigitalLocation(context.Background(), "test-user", location)

		// Verify
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("RemoveDigitalLocation - Success", func(t *testing.T) {
		// Setup
		service := newMockGameDigitalServiceWithDefaults(logger)

		// Execute
		err := service.RemoveDigitalLocation(context.Background(), "test-user", "test-location")

		// Verify
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("RemoveDigitalLocation - Not Found", func(t *testing.T) {
		// Setup
		service := newMockGameDigitalServiceWithDefaults(logger)
		expectedErr := sql.ErrNoRows

		service.dbAdapter = &MockDigitalDbAdapter{
			RemoveDigitalLocationFunc: func(ctx context.Context, userID, locationID string) error {
				return expectedErr
			},
		}

		// Execute
		err := service.RemoveDigitalLocation(context.Background(), "test-user", "test-location")

		// Verify
		if err == nil {
			t.Error("Expected error, got nil")
		}
		if !errors.Is(err, expectedErr) {
			t.Errorf("Expected error %v, got %v", expectedErr, err)
		}
	})

	t.Run("GetUserDigitalLocations - Cache Hit", func(t *testing.T) {
		// Setup
		service := newMockGameDigitalServiceWithDefaults(logger)
		expectedLocations := []models.DigitalLocation{
			{ID: "1", Name: "Location 1"},
			{ID: "2", Name: "Location 2"},
		}

		// Execute
		locations, err := service.GetUserDigitalLocations(context.Background(), "test-user")

		// Verify
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(locations) != len(expectedLocations) {
			t.Errorf("Expected %d locations, got %d", len(expectedLocations), len(locations))
		}
	})

	t.Run("GetUserDigitalLocations - Cache Miss", func(t *testing.T) {
		// Setup
		service := newMockGameDigitalServiceWithDefaults(logger)
		expectedLocations := []models.DigitalLocation{
			{ID: "1", Name: "Location 1"},
			{ID: "2", Name: "Location 2"},
		}

		service.dbAdapter = &MockDigitalDbAdapter{
			GetUserDigitalLocationsFunc: func(ctx context.Context, userID string) ([]models.DigitalLocation, error) {
				return expectedLocations, nil
			},
		}

		// Execute
		locations, err := service.GetUserDigitalLocations(context.Background(), "test-user")

		// Verify
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(locations) != len(expectedLocations) {
			t.Errorf("Expected %d locations, got %d", len(expectedLocations), len(locations))
		}
	})
}

func TestGameDigitalService_GetDigitalLocation(t *testing.T) {
	// Setup
	service := newMockGameDigitalServiceWithDefaults(testutils.NewTestLogger())
	expectedErr := sql.ErrNoRows
	expectedLocation := models.DigitalLocation{}

	service.dbAdapter = &MockDigitalDbAdapter{
		GetDigitalLocationFunc: func(ctx context.Context, userID, locationID string) (models.DigitalLocation, error) {
			return models.DigitalLocation{}, expectedErr
		},
	}

	// Execute
	location, err := service.GetDigitalLocation(context.Background(), "test-user", "test-location")

	// Verify
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if !errors.Is(err, expectedErr) {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}
	if location.ID != "" || location.Name != "" || location.IsActive != false {
		t.Errorf("Expected empty location, got %v", location)
	}

	service.dbAdapter = &MockDigitalDbAdapter{
		GetDigitalLocationFunc: func(ctx context.Context, userID, locationID string) (models.DigitalLocation, error) {
			return expectedLocation, nil
		},
	}

	// Execute
	location, err = service.GetDigitalLocation(context.Background(), "test-user", "test-location")

	// Verify
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if location.ID == "" {
		t.Error("Expected location to be returned")
	}
}

func TestGameDigitalService_AddDigitalLocation(t *testing.T) {
	// Setup
	service := newMockGameDigitalServiceWithDefaults(testutils.NewTestLogger())
	expectedErr := errors.New("test error")
	location := models.DigitalLocation{
		Name:     "Test Location",
		IsActive: true,
	}

	service.dbAdapter = &MockDigitalDbAdapter{
		FindDigitalLocationByNameFunc: func(ctx context.Context, userID string, name string) (models.DigitalLocation, error) {
			return models.DigitalLocation{}, expectedErr
		},
	}

	// Execute
	createdLocation, err := service.AddDigitalLocation(context.Background(), "test-user", location)

	// Verify
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if !errors.Is(err, expectedErr) {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}
	if createdLocation.ID != "" || createdLocation.Name != "" || createdLocation.IsActive != false {
		t.Errorf("Expected empty location, got %v", createdLocation)
	}

	service.dbAdapter = &MockDigitalDbAdapter{
		FindDigitalLocationByNameFunc: func(ctx context.Context, userID string, name string) (models.DigitalLocation, error) {
			return models.DigitalLocation{}, nil
		},
		AddDigitalLocationFunc: func(ctx context.Context, userID string, location models.DigitalLocation) (models.DigitalLocation, error) {
			return location, nil
		},
	}

	// Execute
	createdLocation, err = service.AddDigitalLocation(context.Background(), "test-user", location)

	// Verify
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if createdLocation.Name != location.Name {
		t.Errorf("Expected location name %s, got %s", location.Name, createdLocation.Name)
	}
}

func TestGameDigitalService_UpdateDigitalLocation(t *testing.T) {
	// Setup
	service := newMockGameDigitalServiceWithDefaults(testutils.NewTestLogger())
	expectedErr := sql.ErrNoRows
	location := models.DigitalLocation{
		ID:       "test-location",
		Name:     "Updated Location",
		IsActive: true,
	}

	service.dbAdapter = &MockDigitalDbAdapter{
		GetDigitalLocationFunc: func(ctx context.Context, userID, locationID string) (models.DigitalLocation, error) {
			return models.DigitalLocation{}, expectedErr
		},
	}

	// Execute
	err := service.UpdateDigitalLocation(context.Background(), "test-user", location)

	// Verify
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if !errors.Is(err, expectedErr) {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}

	service.dbAdapter = &MockDigitalDbAdapter{
		GetDigitalLocationFunc: func(ctx context.Context, userID, locationID string) (models.DigitalLocation, error) {
			return models.DigitalLocation{}, nil
		},
		UpdateDigitalLocationFunc: func(ctx context.Context, userID string, location models.DigitalLocation) error {
			return nil
		},
	}

	// Execute
	err = service.UpdateDigitalLocation(context.Background(), "test-user", location)

	// Verify
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestGameDigitalService_RemoveDigitalLocation(t *testing.T) {
	// Setup
	service := newMockGameDigitalServiceWithDefaults(testutils.NewTestLogger())
	expectedErr := sql.ErrNoRows

	service.dbAdapter = &MockDigitalDbAdapter{
		GetDigitalLocationFunc: func(ctx context.Context, userID, locationID string) (models.DigitalLocation, error) {
			return models.DigitalLocation{}, expectedErr
		},
	}

	// Execute
	err := service.RemoveDigitalLocation(context.Background(), "test-user", "test-location")

	// Verify
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if !errors.Is(err, expectedErr) {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}

	service.dbAdapter = &MockDigitalDbAdapter{
		GetDigitalLocationFunc: func(ctx context.Context, userID, locationID string) (models.DigitalLocation, error) {
			return models.DigitalLocation{}, nil
		},
		RemoveDigitalLocationFunc: func(ctx context.Context, userID, locationID string) error {
			return nil
		},
	}

	// Execute
	err = service.RemoveDigitalLocation(context.Background(), "test-user", "test-location")

	// Verify
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}
