export const TOAST_ERROR_MESSAGES = {
  DIGITAL_LOCATION: {
    CREATE: {
      DEFAULT: 'Failed to create digital service',
      PERMISSION: 'You don\'t have permission to create digital services',
      EXISTS: 'A service with these details already exists',
      SERVER: 'Server error occurred while creating service',
      NETWORK: 'Network error. Please check your connection'
    },
    UPDATE: {
      DEFAULT: 'Failed to update digital service',
      PERMISSION: 'You don\'t have permission to update this service',
      NOT_FOUND: 'The service you\'re trying to update doesn\'t exist',
      SERVER: 'Server error occurred while updating service',
      NETWORK: 'Network error. Please check your connection'
    },
    DELETE: {
      DEFAULT: 'Failed to delete digital service',
      PERMISSION: 'You don\'t have permission to delete this service',
      NOT_FOUND: 'The service you\'re trying to delete doesn\'t exist',
      SERVER: 'Server error occurred while deleting service',
      NETWORK: 'Network error. Please check your connection'
    }
  },
  PHYSICAL_LOCATION: {
    CREATE: {
      DEFAULT: 'Failed to create physical location',
      PERMISSION: 'You don\'t have permission to create physical locations',
      EXISTS: 'A location with these details already exists',
      SERVER: 'Server error occurred while creating location',
      NETWORK: 'Network error. Please check your connection'
    },
    UPDATE: {
      DEFAULT: 'Failed to update physical location',
      PERMISSION: 'You don\'t have permission to update this location',
      NOT_FOUND: 'The location you\'re trying to update doesn\'t exist',
      SERVER: 'Server error occurred while updating location',
      NETWORK: 'Network error. Please check your connection'
    },
    DELETE: {
      DEFAULT: 'Failed to delete physical location',
      PERMISSION: 'You don\'t have permission to delete this location',
      NOT_FOUND: 'The location you\'re trying to delete doesn\'t exist',
      SERVER: 'Server error occurred while deleting location',
      NETWORK: 'Network error. Please check your connection'
    }
  },
  SUBLOCATION: {
    CREATE: {
      DEFAULT: 'Failed to create sublocation',
      PERMISSION: 'You don\'t have permission to create sublocations',
      EXISTS: 'A sublocation with these details already exists',
      SERVER: 'Server error occurred while creating sublocation',
      NETWORK: 'Network error. Please check your connection'
    },
    UPDATE: {
      DEFAULT: 'Failed to update sublocation',
      PERMISSION: 'You don\'t have permission to update this sublocation',
      NOT_FOUND: 'The sublocation you\'re trying to update doesn\'t exist',
      SERVER: 'Server error occurred while updating sublocation',
      NETWORK: 'Network error. Please check your connection'
    },
    DELETE: {
      DEFAULT: 'Failed to delete sublocation',
      PERMISSION: 'You don\'t have permission to delete this sublocation',
      NOT_FOUND: 'The sublocation you\'re trying to delete doesn\'t exist',
      SERVER: 'Server error occurred while deleting sublocation',
      NETWORK: 'Network error. Please check your connection'
    }
  },
  GAME: {
    ADD_TO_LIBRARY: {
      DEFAULT: 'Failed to add game to library',
      PERMISSION: 'You don\'t have permission to add games',
      EXISTS: 'Game already exists in your library',
      SERVER: 'Server error occurred while adding game',
      NETWORK: 'Network error. Please check your connection'
    }
  }
} as const;

export type ErrorMessageKey = keyof typeof TOAST_ERROR_MESSAGES;