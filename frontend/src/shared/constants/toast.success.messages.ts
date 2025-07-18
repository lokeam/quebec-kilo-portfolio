export const TOAST_SUCCESS_MESSAGES = {
  DIGITAL_LOCATION: {
    CREATE: 'Digital service created successfully',
    UPDATE: 'Digital service updated successfully',
    DELETE: 'Digital service deleted successfully'
  },
  PHYSICAL_LOCATION: {
    CREATE: 'Physical location created successfully',
    UPDATE: 'Physical location updated successfully',
    DELETE: 'Physical location deleted successfully'
  },
  SUBLOCATION: {
    CREATE: 'Sublocation created successfully',
    UPDATE: 'Sublocation updated successfully',
    DELETE: 'Sublocation deleted successfully'
  },
  GAME: {
    ADD_TO_LIBRARY: 'Game added to your library',
    REMOVE_FROM_LIBRARY: 'Game removed from your library',
    UPDATE: 'Game details updated'
  },
  SPEND_TRACKING: {
    ADD_ITEM: 'Spend item added successfully',
    UPDATE_ITEM: 'Spend item updated successfully',
    DELETE_ITEM: 'Spend item deleted successfully'
  },
  USER: {
    DELETION_REQUESTED: 'Account deletion requested successfully',
    DELETION_CANCELLED: 'Account deletion cancelled successfully'
  }
} as const;

export type ToastSuccessMesssageKey = keyof typeof TOAST_SUCCESS_MESSAGES;
