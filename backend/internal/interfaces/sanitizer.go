package interfaces

type Sanitizer interface {
	SanitizeSearchQuery(query string) (string, error)
}