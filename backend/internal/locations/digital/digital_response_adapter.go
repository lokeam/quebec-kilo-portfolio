package digital

import (
	"github.com/lokeam/qko-beta/internal/types"
)

type DigitalResponseAdapter struct{}

func NewDigitalResponseAdapter() *DigitalResponseAdapter {
	return &DigitalResponseAdapter{}
}

// AdaptToCatalogResponse transforms the digital services catalog into a frontend-friendly format
func (a *DigitalResponseAdapter) AdaptToCatalogResponse(
	catalog []types.DigitalServiceItem,
) []types.DigitalServiceItem {
	return catalog
}