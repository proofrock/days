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

import type { Entry } from './types';

const API_BASE = '/api';

// Get entry for a specific date
export async function getEntry(date: string): Promise<Entry | null> {
  const response = await fetch(`${API_BASE}/entries/${date}`);
  if (response.status === 404) {
    return null;
  }
  if (!response.ok) {
    throw new Error('Failed to fetch entry');
  }
  return response.json();
}

// Save entry (create or update)
export async function saveEntry(entry: Entry): Promise<Entry> {
  const response = await fetch(`${API_BASE}/entries/${entry.date}`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(entry),
  });
  if (!response.ok) {
    throw new Error('Failed to save entry');
  }
  return response.json();
}

// Delete entry
export async function deleteEntry(date: string): Promise<void> {
  const response = await fetch(`${API_BASE}/entries/${date}`, {
    method: 'DELETE',
  });
  if (!response.ok) {
    throw new Error('Failed to delete entry');
  }
}

// Get all entry dates for a specific month
export async function getEntriesByMonth(year: number, month: number): Promise<string[]> {
  const response = await fetch(`${API_BASE}/entries/month/${year}/${month}`);
  if (!response.ok) {
    throw new Error('Failed to fetch entries');
  }
  return response.json();
}

// Get entries summary for a specific month (with working status)
export async function getEntriesSummaryByMonth(
  year: number,
  month: number
): Promise<Array<{ date: string; working: string }>> {
  const response = await fetch(`${API_BASE}/entries/month/${year}/${month}/summary`);
  if (!response.ok) {
    throw new Error('Failed to fetch entries summary');
  }
  return response.json();
}
