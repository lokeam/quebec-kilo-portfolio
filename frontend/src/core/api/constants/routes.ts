export const API_ROUTES = Object.freeze({
  SEARCH: {
    GAMES: '/api/v1/search',
  },
  LOCATIONS: {
    BASE: 'api/v1locations/physical',
    BY_ID: (id: string) => `api/v1/locations/physical/${id}`
  },
  SUBLOCATION: {
    BASE: 'api/v1locations/physical',
    BY_ID: (id: string) => `api/v1/locations/sublocations/${id}`
  },
  DIGITAL: {
    BASE: 'api/v1locations/physical',
    BY_ID: (id: string) => `api/v1/locations/digital/${id}`
  },
});

