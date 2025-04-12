export const API_ROUTES = Object.freeze({
  SEARCH: {
    GAMES: '/api/v1/search',
  },
  LOCATIONS: {
    BASE: '/api/v1/locations',
    BY_ID: (id: string) => `/api/v1/locations/${id}`,
    CREATE: '/api/v1/locations',
  },
  SUBLOCATION: {
    BASE: '/api/v1/sublocations',
    BY_ID: (id: string) => `/api/v1/sublocations/${id}`,
    CREATE: '/api/v1/sublocations',
  },
  DIGITAL: {
    BASE: 'api/v1/locations/physical',
    BY_ID: (id: string) => `api/v1/locations/digital/${id}`
  },
  DIGITAL_SERVICES: {
    BASE: '/api/v1/digital-services',
    BY_ID: (id: string) => `/api/v1/digital-services/${id}`,
    CREATE: '/api/v1/digital-services',
  },
});

