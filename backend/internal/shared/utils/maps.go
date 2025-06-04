package utils

import (
	"fmt"
	"strconv"
	"strings"
)

// BuildGoogleMapsURL takes a latitude and longitude and returns
// a Google Maps search URL pointing to that coordinate.
// Example output: "https://www.google.com/maps/search/?api=1&query=34.410634,132.474688"
func BuildGoogleMapsURL(latitude, longitude float64) string {
	// Format coordinates with 6 decimal places for precision
	query := fmt.Sprintf("%.6f,%.6f", latitude, longitude)

	// Construct the URL directly with the proper format
	return fmt.Sprintf("https://www.google.com/maps/search/?api=1&query=%s", query)
}

// ParseCoordinates takes a coordinate string in the format "latitude, longitude"
// and returns the latitude and longitude as float64 values.
// Returns an error if the string is not in the correct format.
func ParseCoordinates(coords string) (float64, float64, error) {
	parts := strings.Split(strings.TrimSpace(coords), ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid coordinate format: %s", coords)
	}

	lat, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid latitude: %s", parts[0])
	}

	lng, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid longitude: %s", parts[1])
	}

	return lat, lng, nil
}