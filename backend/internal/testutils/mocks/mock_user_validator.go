package mocks

import (
	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/types"
)

type MockUserValidator struct {
	ValidateUserProfileFunc              func(user models.User) (models.User, error)
	ValidateCreateUserRequestFunc        func(request types.CreateUserRequest) (types.CreateUserRequest, error)
	ValidateUpdateUserProfileRequestFunc func(request types.UpdateUserProfileRequest) (types.UpdateUserProfileRequest, error)
}

// ValidateUserProfile mocks the ValidateUserProfile method
func (m *MockUserValidator) ValidateUserProfile(user models.User) (models.User, error) {
	if m.ValidateUserProfileFunc != nil {
		return m.ValidateUserProfileFunc(user)
	}
	return user, nil
}

// ValidateCreateUserRequest mocks the ValidateCreateUserRequest method
func (m *MockUserValidator) ValidateCreateUserRequest(request types.CreateUserRequest) (types.CreateUserRequest, error) {
	if m.ValidateCreateUserRequestFunc != nil {
		return m.ValidateCreateUserRequestFunc(request)
	}
	return request, nil
}

// ValidateUpdateUserProfileRequest mocks the ValidateUpdateUserProfileRequest method
func (m *MockUserValidator) ValidateUpdateUserProfileRequest(request types.UpdateUserProfileRequest) (types.UpdateUserProfileRequest, error) {
	if m.ValidateUpdateUserProfileRequestFunc != nil {
		return m.ValidateUpdateUserProfileRequestFunc(request)
	}
	return request, nil
}