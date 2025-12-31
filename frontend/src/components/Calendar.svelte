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
  import { onMount } from 'svelte';
  import { getEntriesSummaryByMonth } from '../api';

  // Props
  interface Props {
    selectedDate: string;
    onDateSelect: (date: string) => void;
    refreshKey?: number;
  }

  let { selectedDate = $bindable(), onDateSelect, refreshKey = 0 }: Props = $props();

  // State
  let currentYear = $state(new Date().getFullYear());
  let currentMonth = $state(new Date().getMonth() + 1); // 1-12
  let entryData = $state<Map<string, string>>(new Map()); // date -> working status

  // Computed values
  let monthName = $derived(new Date(currentYear, currentMonth - 1).toLocaleString('default', { month: 'long' }));
  let days = $derived(getDaysInMonth(currentYear, currentMonth));

  // Format date as YYYY-MM-DD
  function formatDate(year: number, month: number, day: number): string {
    return `${year}-${String(month).padStart(2, '0')}-${String(day).padStart(2, '0')}`;
  }

  // Get days in month with padding for calendar grid
  function getDaysInMonth(year: number, month: number) {
    const firstDay = new Date(year, month - 1, 1);
    const lastDay = new Date(year, month, 0);
    const daysInMonth = lastDay.getDate();
    const startingDayOfWeek = firstDay.getDay(); // 0 = Sunday

    // Convert to Monday-based week (0 = Monday, 6 = Sunday)
    const mondayBasedStart = startingDayOfWeek === 0 ? 6 : startingDayOfWeek - 1;

    const days: Array<{ day: number; date: string; isToday: boolean; hasEntry: boolean; working: string }> = [];

    // Add padding for days before month starts
    for (let i = 0; i < mondayBasedStart; i++) {
      days.push({ day: 0, date: '', isToday: false, hasEntry: false, working: '' });
    }

    // Add days of the month
    const today = new Date();
    for (let day = 1; day <= daysInMonth; day++) {
      const date = formatDate(year, month, day);
      const isToday = year === today.getFullYear()
        && month === today.getMonth() + 1
        && day === today.getDate();
      const working = entryData.get(date) || '';
      const hasEntry = entryData.has(date);
      days.push({ day, date, isToday, hasEntry, working });
    }

    return days;
  }

  // Load entries for current month
  async function loadEntries() {
    try {
      const summaries = await getEntriesSummaryByMonth(currentYear, currentMonth);
      entryData = new Map(summaries.map(s => [s.date, s.working]));
    } catch (error) {
      console.error('Failed to load entries:', error);
    }
  }

  // Navigate to previous month
  function previousMonth() {
    if (currentMonth === 1) {
      currentMonth = 12;
      currentYear--;
    } else {
      currentMonth--;
    }
  }

  // Navigate to next month
  function nextMonth() {
    if (currentMonth === 12) {
      currentMonth = 1;
      currentYear++;
    } else {
      currentMonth++;
    }
  }

  // Handle day click
  function selectDay(date: string) {
    if (date) {
      selectedDate = date;
      onDateSelect(date);
    }
  }

  // Load entries when month changes or refreshKey changes
  $effect(() => {
    const year = currentYear;
    const month = currentMonth;
    const key = refreshKey;
    loadEntries();
  });

  // Initialize with today's date selected
  onMount(() => {
    const today = new Date();
    selectedDate = formatDate(today.getFullYear(), today.getMonth() + 1, today.getDate());
    onDateSelect(selectedDate);
  });
</script>

<div class="calendar-container">
  <!-- Month navigation -->
  <div class="d-flex justify-content-between align-items-center mb-3">
    <button class="btn btn-outline-secondary" onclick={previousMonth}>
      &larr;
    </button>
    <h4 class="mb-0">{monthName} {currentYear}</h4>
    <button class="btn btn-outline-secondary" onclick={nextMonth}>
      &rarr;
    </button>
  </div>

  <!-- Calendar grid -->
  <div class="calendar-grid">
    <!-- Day headers -->
    {#each ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'] as dayName}
      <div class="calendar-day-header">{dayName}</div>
    {/each}

    <!-- Days -->
    {#each days as { day, date, isToday, hasEntry, working }}
      {#if day === 0}
        <div class="calendar-day empty"></div>
      {:else}
        <button
          class="calendar-day"
          class:today={isToday}
          class:has-entry={hasEntry}
          class:worked={working === 'yes'}
          class:worked-partial={working === 'partial'}
          class:not-worked={working === 'no'}
          class:selected={date === selectedDate}
          onclick={() => selectDay(date)}
        >
          {day}
        </button>
      {/if}
    {/each}
  </div>
</div>

<style>
  .calendar-container {
    max-width: 600px;
    margin: 0 auto;
  }

  .calendar-grid {
    display: grid;
    grid-template-columns: repeat(7, 1fr);
    gap: 4px;
  }

  .calendar-day-header {
    text-align: center;
    font-weight: bold;
    padding: 8px;
    background-color: #f8f9fa;
    border-radius: 4px;
  }

  .calendar-day {
    aspect-ratio: 1;
    border: 1px solid #dee2e6;
    background-color: white;
    border-radius: 4px;
    cursor: pointer;
    transition: all 0.2s;
    font-size: 14px;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .calendar-day.empty {
    border: none;
    cursor: default;
  }

  .calendar-day:not(.empty):hover {
    background-color: #e9ecef;
  }

  .calendar-day.today {
    border-color: #0d6efd;
    font-weight: bold;
  }

  .calendar-day.has-entry {
    background-color: #9fd3e5;
  }

  .calendar-day.worked {
    background-color: #d1e7dd;
    border-color: #198754;
  }

  .calendar-day.worked-partial {
    background-color: #fff3cd;
    border-color: #ffc107;
  }

  .calendar-day.not-worked {
    background-color: #f8d7da;
    border-color: #dc3545;
  }

  .calendar-day.selected {
    background-color: #0d6efd;
    color: white;
    border-color: #0d6efd;
  }

  .calendar-day.selected.has-entry {
    background-color: #0a58ca;
  }
</style>
