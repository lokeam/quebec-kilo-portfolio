package media_storage

import (
	"github.com/lokeam/qko-beta/internal/analytics"
)

// MediaStorageResponse represents the complete response for the media storage page
// It reuses the StorageStats type from analytics since it already contains
// all the necessary information about storage locations
type MediaStorageResponse struct {
	// StorageStats contains information about both physical and digital storage locations
    // including counts, item counts, and subscription details
		StorageStats *analytics.StorageStats `json:"storage_stats"`
}