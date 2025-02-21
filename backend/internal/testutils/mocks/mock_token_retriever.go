package mocks

import (
	"context"

	"github.com/lokeam/qko-beta/internal/interfaces"
)

type MockTwitchTokenRetriever struct{}

func (m *MockTwitchTokenRetriever) GetToken(
    ctx context.Context,
    clientID,
    clientSecret,
    authURL string,
    logger interfaces.Logger,
) (string, error) {
	return "mock-valid-access-token", nil
}
