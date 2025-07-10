package mocks

import (
	"context"

	"github.com/lokeam/qko-beta/internal/models"
)

type MockUserDbAdapter struct {
	GetSingleUserFunc        func(ctx context.Context, userID string) (models.User, error)
	CreateUserFunc           func(ctx context.Context, user models.User) (models.User, error)
	UpdateUserProfileFunc    func(ctx context.Context, userID string, firstName, lastName string) (models.User, error)
	HasCompleteProfileFunc   func(ctx context.Context, userID string) (bool, error)
	GetSingleUserByEmailFunc func(ctx context.Context, email string) (models.User, error)
}

// GetSingleUser mocks the GetSingleUser method
func (m *MockUserDbAdapter) GetSingleUser(ctx context.Context, userID string) (models.User, error) {
	if m.GetSingleUserFunc != nil {
		return m.GetSingleUserFunc(ctx, userID)
	}
	return models.User{}, nil
}

// CreateUser mocks the CreateUser method
func (m *MockUserDbAdapter) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	if m.CreateUserFunc != nil {
		return m.CreateUserFunc(ctx, user)
	}
	return user, nil
}

// UpdateUserProfile mocks the UpdateUserProfile method
func (m *MockUserDbAdapter) UpdateUserProfile(ctx context.Context, userID string, firstName, lastName string) (models.User, error) {
	if m.UpdateUserProfileFunc != nil {
		return m.UpdateUserProfileFunc(ctx, userID, firstName, lastName)
	}
	return models.User{}, nil
}

// HasCompleteProfile mocks the HasCompleteProfile method
func (m *MockUserDbAdapter) HasCompleteProfile(ctx context.Context, userID string) (bool, error) {
	if m.HasCompleteProfileFunc != nil {
		return m.HasCompleteProfileFunc(ctx, userID)
	}
	return true, nil
}

// GetSingleUserByEmail mocks the GetSingleUserByEmail method
func (m *MockUserDbAdapter) GetSingleUserByEmail(ctx context.Context, email string) (models.User, error) {
	if m.GetSingleUserByEmailFunc != nil {
		return m.GetSingleUserByEmailFunc(ctx, email)
	}
	return models.User{}, nil
}