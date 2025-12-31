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
  import type { Entry, ViewMode } from '../types';
  import { FIELD_IDS } from '../types';
  import { getEntry, saveEntry } from '../api';
  import StarRating from './StarRating.svelte';

  // Props
  interface Props {
    date: string;
    mode: ViewMode;
    onClose: () => void;
    onSave?: () => void;
  }

  let { date, mode, onClose, onSave }: Props = $props();

  // State
  let loading = $state(true);
  let saving = $state(false);

  // Form field values (separate state for easier manipulation)
  let positionLon = $state('');
  let positionLat = $state('');
  let positionName = $state('');
  let rating = $state(0);
  let general = $state('');
  let working = $state(''); // '', 'yes', 'partial', 'no'
  let mood = $state(0);
  let moodTxt = $state('');
  let lunch = $state('');
  let dinner = $state('');
  let tv = $state('');
  let sleep = $state(0);
  let sleepTxt = $state('');

  // Computed
  let isViewMode = $derived(mode === 'view');
  let hasPosition = $derived(positionLon !== '' && positionLat !== '');
  let googleMapsUrl = $derived(
    hasPosition
      ? `https://www.google.com/maps?q=${positionLat},${positionLon}`
      : ''
  );

  // Load entry data
  async function loadEntry() {
    loading = true;
    try {
      const data = await getEntry(date);
      if (data) {
        // Populate form fields
        positionLon = data.fields[FIELD_IDS.POSITION_LON] || '';
        positionLat = data.fields[FIELD_IDS.POSITION_LAT] || '';
        positionName = data.fields[FIELD_IDS.POSITION_NAME] || '';
        rating = parseInt(data.fields[FIELD_IDS.RATING] || '0');
        general = data.fields[FIELD_IDS.GENERAL] || '';
        working = data.fields[FIELD_IDS.WORKING] || '';
        mood = parseInt(data.fields[FIELD_IDS.MOOD] || '0');
        moodTxt = data.fields[FIELD_IDS.MOOD_TXT] || '';
        lunch = data.fields[FIELD_IDS.LUNCH] || '';
        dinner = data.fields[FIELD_IDS.DINNER] || '';
        tv = data.fields[FIELD_IDS.TV] || '';
        sleep = parseInt(data.fields[FIELD_IDS.SLEEP] || '0');
        sleepTxt = data.fields[FIELD_IDS.SLEEP_TXT] || '';
      }
    } catch (error) {
      console.error('Failed to load entry:', error);
    } finally {
      loading = false;
    }
  }

  // Get current position
  async function getCurrentPosition() {
    if (!navigator.geolocation) {
      alert('Geolocation is not supported by your browser');
      return;
    }

    navigator.geolocation.getCurrentPosition(
      (position) => {
        positionLat = position.coords.latitude.toString();
        positionLon = position.coords.longitude.toString();
      },
      (error) => {
        alert(`Failed to get location: ${error.message}`);
      }
    );
  }

  // Reverse geocode the current coordinates
  async function reverseGeocode() {
    if (!positionLat || !positionLon) {
      alert('Please enter or get coordinates first');
      return;
    }

    try {
      const lat = parseFloat(positionLat);
      const lon = parseFloat(positionLon);

      if (isNaN(lat) || isNaN(lon)) {
        alert('Invalid coordinates');
        return;
      }

      const response = await fetch(
        `https://nominatim.openstreetmap.org/reverse?format=json&lat=${lat}&lon=${lon}&zoom=14&addressdetails=1`,
        {
          headers: {
            'User-Agent': 'Days Journal App'
          }
        }
      );

      if (response.ok) {
        const data = await response.json();
        // Build a simple location name from the address
        const parts = [];
        if (data.address.village || data.address.town || data.address.city) {
          parts.push(data.address.village || data.address.town || data.address.city);
        }
        if (data.address.country) {
          parts.push(data.address.country);
        }
        positionName = parts.join(', ') || data.display_name;
      } else {
        alert('Failed to get location name');
      }
    } catch (error) {
      console.error('Reverse geocoding failed:', error);
      alert('Failed to get location name');
    }
  }

  // Save entry
  async function handleSave() {
    saving = true;
    try {
      // Build fields object
      const fields: Record<string, string> = {
        [FIELD_IDS.POSITION_LON]: positionLon,
        [FIELD_IDS.POSITION_LAT]: positionLat,
        [FIELD_IDS.POSITION_NAME]: positionName,
        [FIELD_IDS.RATING]: rating.toString(),
        [FIELD_IDS.GENERAL]: general,
        [FIELD_IDS.WORKING]: working,
        [FIELD_IDS.MOOD]: mood.toString(),
        [FIELD_IDS.MOOD_TXT]: moodTxt,
        [FIELD_IDS.LUNCH]: lunch,
        [FIELD_IDS.DINNER]: dinner,
        [FIELD_IDS.TV]: tv,
        [FIELD_IDS.SLEEP]: sleep.toString(),
        [FIELD_IDS.SLEEP_TXT]: sleepTxt,
      };

      const entryToSave: Entry = {
        date,
        timestamp: new Date().toISOString(),
        fields,
      };

      await saveEntry(entryToSave);
      onSave?.();
      onClose();
    } catch (error) {
      console.error('Failed to save entry:', error);
      alert('Failed to save entry');
    } finally {
      saving = false;
    }
  }

  // Load entry on mount
  onMount(() => {
    loadEntry();
  });
</script>

<div class="entry-form-container">
  <div class="card">
    <div class="card-header">
      <h5 class="mb-0">{mode === 'view' ? 'View' : 'Edit'} Entry - {date}</h5>
    </div>
    <div class="card-body">
      {#if loading}
        <div class="text-center">
          <div class="spinner-border" role="status">
            <span class="visually-hidden">Loading...</span>
          </div>
        </div>
      {:else}
        <form onsubmit={(e) => { e.preventDefault(); handleSave(); }}>
          <!-- Buttons -->
          <div class="d-flex gap-2 justify-content-center mb-4">
            <button type="button" class="btn btn-secondary btn-lg" onclick={onClose}>
              Cancel
            </button>
            {#if !isViewMode}
              <button type="submit" class="btn btn-primary btn-lg" disabled={saving}>
                {saving ? 'Saving...' : 'Save'}
              </button>
            {/if}
          </div>

          <!-- Position (Latitude/Longitude) -->
          <div class="mb-3">
            <div class="fw-bold mb-2">Position</div>
            <div class="row g-2 mb-2">
              <div class="col-md-5">
                <input
                  type="text"
                  class="form-control"
                  placeholder="Latitude"
                  bind:value={positionLat}
                  disabled={isViewMode}
                />
              </div>
              <div class="col-md-5">
                <input
                  type="text"
                  class="form-control"
                  placeholder="Longitude"
                  bind:value={positionLon}
                  disabled={isViewMode}
                />
              </div>
              {#if !isViewMode}
                <div class="col-md-1">
                  <button type="button" class="btn btn-outline-primary w-100" onclick={getCurrentPosition} title="Get current position">
                    📍
                  </button>
                </div>
                <div class="col-md-1">
                  <button type="button" class="btn btn-outline-secondary w-100" onclick={reverseGeocode} title="Get location name" disabled={!hasPosition}>
                    🏷️
                  </button>
                </div>
              {/if}
            </div>
            {#if hasPosition}
              <div class="mt-2">
                {#if positionName}
                  <div class="text-muted small mb-1">
                    <strong>📍 {positionName}</strong>
                  </div>
                {/if}
                <a href={googleMapsUrl} target="_blank" rel="noopener noreferrer" class="btn btn-link btn-sm ps-0">
                  View on Google Maps
                </a>
              </div>
            {/if}
          </div>

          <!-- Rating -->
          <div class="mb-3">
            <div class="fw-bold mb-2">Rating</div>
            <StarRating bind:value={rating} disabled={isViewMode} />
          </div>

          <!-- General -->
          <div class="mb-3">
            <label for="general-field" class="form-label fw-bold">General</label>
            <textarea
              id="general-field"
              class="form-control"
              rows="12"
              bind:value={general}
              disabled={isViewMode}
            ></textarea>
          </div>

          <!-- Working -->
          <div class="mb-3">
            <div class="fw-bold mb-2">Working</div>
            <div class="btn-group" role="group">
              <button
                type="button"
                class="btn"
                class:btn-success={working === 'yes'}
                class:btn-outline-secondary={working !== 'yes'}
                onclick={() => working = working === 'yes' ? '' : 'yes'}
                disabled={isViewMode}
              >
                Worked
              </button>
              <button
                type="button"
                class="btn"
                class:btn-warning={working === 'partial'}
                class:btn-outline-secondary={working !== 'partial'}
                onclick={() => working = working === 'partial' ? '' : 'partial'}
                disabled={isViewMode}
              >
                Partial
              </button>
              <button
                type="button"
                class="btn"
                class:btn-danger={working === 'no'}
                class:btn-outline-secondary={working !== 'no'}
                onclick={() => working = working === 'no' ? '' : 'no'}
                disabled={isViewMode}
              >
                Not worked
              </button>
            </div>
          </div>

          <!-- Mood -->
          <div class="mb-3">
            <div class="fw-bold mb-2">Mood</div>
            <StarRating bind:value={mood} disabled={isViewMode} />
          </div>

          <!-- Mood Text -->
          <div class="mb-3">
            <label for="mood-txt-field" class="form-label fw-bold">Mood Description</label>
            <input
              id="mood-txt-field"
              type="text"
              class="form-control"
              bind:value={moodTxt}
              disabled={isViewMode || mood === 0}
            />
          </div>

          <!-- Lunch -->
          <div class="mb-3">
            <label for="lunch-field" class="form-label fw-bold">Lunch</label>
            <textarea
              id="lunch-field"
              class="form-control"
              rows="2"
              bind:value={lunch}
              disabled={isViewMode}
            ></textarea>
          </div>

          <!-- Dinner -->
          <div class="mb-3">
            <label for="dinner-field" class="form-label fw-bold">Dinner</label>
            <textarea
              id="dinner-field"
              class="form-control"
              rows="2"
              bind:value={dinner}
              disabled={isViewMode}
            ></textarea>
          </div>

          <!-- TV -->
          <div class="mb-3">
            <label for="tv-field" class="form-label fw-bold">TV</label>
            <textarea
              id="tv-field"
              class="form-control"
              rows="2"
              bind:value={tv}
              disabled={isViewMode}
            ></textarea>
          </div>

          <!-- Sleep -->
          <div class="mb-3">
            <div class="fw-bold mb-2">Sleep</div>
            <StarRating bind:value={sleep} disabled={isViewMode} />
          </div>

          <!-- Sleep Text -->
          <div class="mb-3">
            <label for="sleep-txt-field" class="form-label fw-bold">Sleep Description</label>
            <input
              id="sleep-txt-field"
              type="text"
              class="form-control"
              bind:value={sleepTxt}
              disabled={isViewMode || sleep === 0}
            />
          </div>
        </form>
      {/if}
    </div>
  </div>
</div>

<style>
  .entry-form-container {
    max-width: 800px;
    margin: 0 auto;
  }
</style>
