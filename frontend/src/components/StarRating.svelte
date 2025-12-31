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
