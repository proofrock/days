// Field IDs for journal entries
export const FIELD_IDS = {
  POSITION_LON: 'POSITION_LON',
  POSITION_LAT: 'POSITION_LAT',
  POSITION_NAME: 'POSITION_NAME',
  RATING: 'RATING',
  GENERAL: 'GENERAL',
  MOOD: 'MOOD',
  MOOD_TXT: 'MOOD_TXT',
  LUNCH: 'LUNCH',
  DINNER: 'DINNER',
  TV: 'TV',
  SLEEP: 'SLEEP',
  SLEEP_TXT: 'SLEEP_TXT',
} as const;

// Journal entry type
export interface Entry {
  date: string; // YYYY-MM-DD format
  timestamp: string; // ISO8601 timestamp
  fields: Record<string, string>; // field_id -> value mapping
}

// View mode for entry form
export type ViewMode = 'view' | 'edit';
