<script>
  export let tabs = [];
  export let active = 0;
</script>

<style>
  .active {
    font-weight: bold;
  }
</style>

<nav>
  {#each tabs as tab, i}
    <button class:active={active === i} on:click={() => active = i}>
      {tab.title}
    </button>
  {/each}
</nav>

<svelte:component
  this={tabs[active].component}
  title={tabs[active].title}
  {...tabs[active].props}
  />
