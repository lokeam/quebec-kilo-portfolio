export const API_ROUTES = Object.freeze({
  SEARCH: {
    GAMES: '/v1/search',
  },
  LOCATIONS: {
    BASE: '/v1/locations/physical',
    BY_ID: (id: string) => `/v1/locations/physical/${id}`,
    CREATE: '/v1/locations/physical',
  },
  SUBLOCATION: {
    BASE: '/v1/locations/sublocations',
    BY_ID: (id: string) => `/v1/locations/sublocations/${id}`,
    CREATE: '/v1/locations/sublocations',
  },
  DIGITAL: {
    BASE: '/v1/locations/digital',
    BY_ID: (id: string) => `/v1/locations/digital/${id}`
  },
  DIGITAL_SERVICES: {
    BASE: '/v1/digital-services',
    BY_ID: (id: string) => `/v1/digital-services/${id}`,
    CREATE: '/v1/digital-services',
  },
});

