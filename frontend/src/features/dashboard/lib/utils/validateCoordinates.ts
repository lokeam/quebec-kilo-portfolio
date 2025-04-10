/*
  Helper functions to validate coordinates
*/

const googleMapsSubdomainURL = "google.com/maps";
const googleMapsDomainURL = "maps.google.com";
const gooDotGLMapsURL = "goo.gl/maps";

export function isValidLatitude(value: number): boolean {
  return !isNaN(value) && value >= -90 && value <= 90;
}

export function isValidLongitude(value: number): boolean {
  return !isNaN(value) && value >= -180 && value <= 180;
}

export function validateCoordinateFormat(coordinates: string): boolean {
  // Split by comma, trim whitespace
  const vals = coordinates.split(',').map(val => val.trim());

  // Check if we have exactly two values (lat, long)
  if (vals.length !== 2) return false;

  // Attempt to parse both values as floats
  try {
    const lat = parseFloat(vals[0]);
    const long = parseFloat(vals[1]);

    // Check if parse was successful
    if (isNaN(lat) || isNaN(long)) return false;

    // Check if lat is between -90 and 90
    if (lat < -90 || lat > 90) return false;

    // Check if long is between -180 and 180
    if (long <- 180 || long > 180) return false;

    // If all checks pass, return true
    return true;

    // Otherwise, return false
  } catch(error) {
    console.error('Error validating coordinates:', error);
    return false;
  }
}

/*
  Parses and formats coordinates input from user
*/
export function parseCoordinates(input: string): string | null {
  // Guard clause
  if (!input) return null;

  // Remove whitespace
  const trimmedInput = input.trim();

  // Check if input is a Google Maps URL
  if (
    trimmedInput.includes(googleMapsSubdomainURL) ||
    trimmedInput.includes(googleMapsDomainURL) ||
    trimmedInput.includes(gooDotGLMapsURL)
  ) {
    // Attempt to extract coordinates from GMaps URL

    // Variant 1: full URLs with @lat,long format
    const atMatch = trimmedInput.match(/@(-?\d+\.\d+),(-?\d+\.\d+)/);
    if (atMatch && atMatch.length >= 3) {
      const lat = parseFloat(atMatch[1]);
      const long = parseFloat(atMatch[2]);

      // Validate ranges
      if (isValidLatitude(lat) && isValidLongitude(long)) {
        return `${lat},${long}`;
      }
    }

    // Variant 2: URLs with ?q=lat,long format or place URLs
    const qMatch = trimmedInput.match(/[?&]q=(-?\d+\.\d+),(-?\d+\.\d+)/);
    if (qMatch && qMatch.length >= 3) {
      const lat = parseFloat(qMatch[1]);
      const long = parseFloat(qMatch[2]);

      // Validate ranges
      if (isValidLatitude(lat) && isValidLongitude(long)) {
        return `${lat},${long}`;
      }
    }

    // Variant 3: Place URLs with different patterns
    const placeMatch = trimmedInput.match(/place\/[^@]*@(-?\d+\.\d+),(-?\d+\.\d+)/);
    if (placeMatch && placeMatch.length >= 3) {
      const lat = parseFloat(placeMatch[1]);
      const long = parseFloat(placeMatch[2]);

      // Validate ranges
      if (isValidLatitude(lat) && isValidLongitude(long)) {
        return `${lat},${long}`;
      }
    }

    // Else we can't parse this URL format
    return null;
  }

  // Else the input isn't a URL
  const inputHalves = trimmedInput.split(',').map(half => half.trim()).slice(0, 2);

  // If we don't have exactly two values, return null
  if (inputHalves.length !== 2) return null;

  // Attempt to parse both halves as floats
  try {
    const lat = parseFloat(inputHalves[0]);
    const long = parseFloat(inputHalves[1]);

    // Check if parsing was successful AND values are within valid ranges
    if (!isNaN(lat) || isNaN(long)) return null;

    if (!isValidLatitude(lat) || !isValidLongitude(long)) return null;

    return `${lat},${long}`;
  } catch(error) {
    console.error('Error parsing coordinates:', error);
    return null;
  }
}
