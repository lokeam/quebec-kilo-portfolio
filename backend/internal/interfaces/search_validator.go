package interfaces

import "github.com/lokeam/qko-beta/internal/search/searchdef"

type SearchValidator interface {
	ValidateQuery(query searchdef.SearchQuery) error
}
