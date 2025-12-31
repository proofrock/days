// Copyright 2025
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Field IDs for journal entries
export const FIELD_IDS = {
  POSITION_LON: 'POSITION_LON',
  POSITION_LAT: 'POSITION_LAT',
  POSITION_NAME: 'POSITION_NAME',
  RATING: 'RATING',
  GENERAL: 'GENERAL',
  WORKING: 'WORKING',
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
