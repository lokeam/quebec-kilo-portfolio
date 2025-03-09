package mocks

// NOTE: Method to MockSanitizer (different from SimpleMockSanitizer)
func (ms *MockSanitizer) SanitizeString(s string) (string, error) {
    return ms.SanitizeFunc(s)
}
