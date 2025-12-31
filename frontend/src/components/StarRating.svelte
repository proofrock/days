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
  // Props
  interface Props {
    value: number; // 0-10, where 0 means no rating
    disabled?: boolean;
    onChange?: (value: number) => void;
  }

  let { value = $bindable(0), disabled = false, onChange }: Props = $props();

  // Handle star click
  function handleClick(rating: number) {
    if (disabled) return;
    value = rating;
    onChange?.(rating);
  }

  // Determine if a star should be filled
  function isFilled(starIndex: number): boolean {
    return value >= starIndex;
  }
</script>

<div class="star-rating">
  {#each Array(10) as _, i}
    <button
      type="button"
      class="star"
      class:filled={isFilled(i + 1)}
      class:disabled={disabled}
      onclick={() => handleClick(i + 1)}
      {disabled}
    >
      {isFilled(i + 1) ? '★' : '☆'}
    </button>
  {/each}
  <span class="rating-value ms-2">
    {value > 0 ? value : 'N.a.'}
  </span>
  {#if value > 0}
    <button
      type="button"
      class="btn btn-sm btn-outline-secondary ms-2"
      onclick={() => handleClick(0)}
      {disabled}
    >
      Clear
    </button>
  {/if}
</div>

<style>
  .star-rating {
    display: flex;
    align-items: center;
    gap: 4px;
  }

  .star {
    background: none;
    border: none;
    font-size: 24px;
    color: #ffc107;
    cursor: pointer;
    padding: 0;
    line-height: 1;
    transition: transform 0.1s;
  }

  .star:not(.disabled):hover {
    transform: scale(1.2);
  }

  .star.disabled {
    cursor: default;
    opacity: 0.6;
  }

  .star:not(.filled) {
    color: #dee2e6;
  }

  .rating-value {
    font-weight: 500;
    min-width: 2.5rem;
    text-align: center;
  }
</style>
