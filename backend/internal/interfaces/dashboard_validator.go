package interfaces

type DashboardValidator interface {
	ValidateUserID(userID string) error
}
