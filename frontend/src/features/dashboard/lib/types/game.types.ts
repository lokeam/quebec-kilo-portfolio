/**
 * Types for game data
 * Aligns with the backend models for games
 */

/**
 * Game model - represents a video game in the user's library
 */
export interface Game {
  id: number;
  name: string;
  summary?: string;
  cover_id?: string;
  cover_url?: string;
  first_release_date?: string;
  rating?: number;
  platform_names?: string[];
  genre_names?: string[];
  theme_names?: string[];
}