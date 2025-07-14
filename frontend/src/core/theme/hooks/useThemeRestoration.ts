/**
 * Synchronously restore theme data from backup before theme store initialization
 * This must run BEFORE any component that uses the theme store renders
 */
export function restoreThemeDataSync() {
  console.log('🔍 Checking for theme backup...');

  // Check for backup data and restore it synchronously
  const backupData = sessionStorage.getItem('qko-theme-backup');
  console.log('🔍 Backup data found:', !!backupData);

  if (backupData) {
    localStorage.setItem('qko-theme-storage', backupData);
    sessionStorage.removeItem('qko-theme-backup');
    return true;
  }

  return false;
}

// Execute immediately when this module is loaded
restoreThemeDataSync();