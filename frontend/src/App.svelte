<!--
Copyright 2025

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
-->

<script lang="ts">
  import type { ViewMode } from './types';
  import { getEntry, deleteEntry as deleteEntryAPI } from './api';
  import Calendar from './components/Calendar.svelte';
  import EntryForm from './components/EntryForm.svelte';

  // State
  let selectedDate = $state('');
  let hasEntry = $state(false);
  let showForm = $state(false);
  let formMode = $state<ViewMode>('edit');
  let calendarRefreshKey = $state(0);

  // Check if selected date has an entry
  async function checkEntry(date: string) {
    try {
      const entry = await getEntry(date);
      hasEntry = entry !== null;
    } catch (error) {
      console.error('Failed to check entry:', error);
      hasEntry = false;
    }
  }

  // Handle date selection from calendar
  function handleDateSelect(date: string) {
    selectedDate = date;
    checkEntry(date);
    showForm = false;
  }

  // Show entry form in view mode
  function viewEntry() {
    formMode = 'view';
    showForm = true;
  }

  // Show entry form in edit mode
  function editEntry() {
    formMode = 'edit';
    showForm = true;
  }

  // Close form and return to calendar
  function closeForm() {
    showForm = false;
    checkEntry(selectedDate); // Refresh entry status
  }

  // Handle save callback
  function handleSave() {
    checkEntry(selectedDate);
    // Refresh calendar to update entry markers
    calendarRefreshKey++;
  }

  // Delete entry with confirmation
  async function deleteEntry() {
    if (!confirm(`Are you sure you want to delete the entry for ${selectedDate}?`)) {
      return;
    }

    try {
      await deleteEntryAPI(selectedDate);
      hasEntry = false;
      // Trigger calendar refresh
      calendarRefreshKey++;
    } catch (error) {
      console.error('Failed to delete entry:', error);
      alert('Failed to delete entry');
    }
  }
</script>

<div class="container py-4">
  <header class="mb-4">
    <h1 class="text-center">Days - Journal</h1>
  </header>

  <main>
    {#if !showForm}
      <!-- Action buttons -->
      <div class="d-flex gap-2 justify-content-center mb-4">
        <button
          class="btn btn-primary btn-lg"
          onclick={viewEntry}
          disabled={!hasEntry}
        >
          View
        </button>
        <button
          class="btn btn-success btn-lg"
          onclick={editEntry}
        >
          Edit
        </button>
        <button
          class="btn btn-danger btn-lg"
          onclick={deleteEntry}
          disabled={!hasEntry}
        >
          Delete
        </button>
      </div>

      <!-- Calendar View -->
      <Calendar bind:selectedDate onDateSelect={handleDateSelect} refreshKey={calendarRefreshKey} />
    {:else}
      <!-- Entry Form -->
      <EntryForm
        date={selectedDate}
        mode={formMode}
        onClose={closeForm}
        onSave={handleSave}
      />
    {/if}
  </main>
</div>

<style>
  :global(body) {
    background-color: #f8f9fa;
  }

  header h1 {
    color: #212529;
    font-weight: 300;
  }
</style>
