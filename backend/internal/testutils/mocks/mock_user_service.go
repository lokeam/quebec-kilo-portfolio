package mocks

import (
	"context"

	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/types"
)

type MockUserService struct {
	GetSingleUserFunc      func(ctx context.Context, userID string) (models.User, error)
	CreateUserFunc         func(ctx context.Context, req types.CreateUserRequest) (models.User, error)
	UpdateUserProfileFunc  func(ctx context.Context, userID string, req types.UpdateUserProfileRequest) (models.User, error)
	HasCompleteProfileFunc func(ctx context.Context, userID string) (bool, error)
}

// GetSingleUser mocks the GetSingleUser method
func (m *MockUserService) GetSingleUser(ctx context.Context, userID string) (models.User, error) {
	if m.GetSingleUserFunc != nil {
		return m.GetSingleUserFunc(ctx, userID)
	}
	return models.User{}, nil
}

// CreateUser mocks the CreateUser method
func (m *MockUserService) CreateUser(ctx context.Context, req types.CreateUserRequest) (models.User, error) {
	if m.CreateUserFunc != nil {
		return m.CreateUserFunc(ctx, req)
	}
	return models.User{}, nil
}

// UpdateUserProfile mocks the UpdateUserProfile method
func (m *MockUserService) UpdateUserProfile(ctx context.Context, userID string, req types.UpdateUserProfileRequest) (models.User, error) {
	if m.UpdateUserProfileFunc != nil {
		return m.UpdateUserProfileFunc(ctx, userID, req)
	}
	return models.User{}, nil
}

// HasCompleteProfile mocks the HasCompleteProfile method
func (m *MockUserService) HasCompleteProfile(ctx context.Context, userID string) (bool, error) {
	if m.HasCompleteProfileFunc != nil {
		return m.HasCompleteProfileFunc(ctx, userID)
	}
	return true, nil
}