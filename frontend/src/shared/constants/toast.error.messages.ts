export const TOAST_ERROR_MESSAGES = {
  DIGITAL_LOCATION: {
    CREATE: {
      DEFAULT: 'Failed to create digital service',
      PERMISSION: 'You don\'t have permission to create digital services',
      EXISTS: 'A service with these details already exists',
      SERVER: 'Server error occurred while creating service',
    },
    UPDATE: {
      DEFAULT: 'Failed to update digital service',
      PERMISSION: 'You don\'t have permission to update this service',
      NOT_FOUND: 'The service you\'re trying to update doesn\'t exist',
      SERVER: 'Server error occurred while updating service',
    },
    DELETE: {
      DEFAULT: 'Failed to delete digital service',
      PERMISSION: 'You don\'t have permission to delete this service',
      NOT_FOUND: 'The service you\'re trying to delete doesn\'t exist',
      SERVER: 'Server error occurred while deleting service',
    }
  },
  PHYSICAL_LOCATION: {
    CREATE: {
      DEFAULT: 'Failed to create physical location',
      PERMISSION: 'You don\'t have permission to create physical locations',
      EXISTS: 'A location with these details already exists',
      SERVER: 'Server error occurred while creating location',
    },
    UPDATE: {
      DEFAULT: 'Failed to update physical location',
      PERMISSION: 'You don\'t have permission to update this location',
      NOT_FOUND: 'The location you\'re trying to update doesn\'t exist',
      SERVER: 'Server error occurred while updating location',
    },
    DELETE: {
      DEFAULT: 'Failed to delete physical location',
      PERMISSION: 'You don\'t have permission to delete this location',
      NOT_FOUND: 'The location you\'re trying to delete doesn\'t exist',
      SERVER: 'Server error occurred while deleting location',
    }
  },
  SUBLOCATION: {
    CREATE: {
      DEFAULT: 'Failed to create sublocation',
      PERMISSION: 'You don\'t have permission to create sublocations',
      EXISTS: 'A sublocation with these details already exists',
      SERVER: 'Server error occurred while creating sublocation',
    },
    UPDATE: {
      DEFAULT: 'Failed to update sublocation',
      PERMISSION: 'You don\'t have permission to update this sublocation',
      NOT_FOUND: 'The sublocation you\'re trying to update doesn\'t exist',
      SERVER: 'Server error occurred while updating sublocation',
    },
    DELETE: {
      DEFAULT: 'Failed to delete sublocation',
      PERMISSION: 'You don\'t have permission to delete this sublocation',
      NOT_FOUND: 'The sublocation you\'re trying to delete doesn\'t exist',
      SERVER: 'Server error occurred while deleting sublocation',
    }
  },
  GAME: {
    ADD_TO_LIBRARY: {
      DEFAULT: 'Failed to add game to library',
      PERMISSION: 'You don\'t have permission to add games',
      EXISTS: 'Game already exists in your library',
      SERVER: 'Server error occurred while adding game',
    }
  }
} as const;

export type ErrorMessageKey = keyof typeof TOAST_ERROR_MESSAGES;