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
