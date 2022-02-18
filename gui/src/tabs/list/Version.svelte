<script>
  export let version = {};

  $: games = version?.games || {};
  $: count = Object.values(games).filter(game => game.isInstalled).length;
</script>

<style lang="scss">
  @import '../../main.scss';

  h2 {
    // font-weight: 500;
    font-size: 18px;

    .count {
      color: $foreground-dark;
    }
  }
</style>

{#if count > 0}
  <h2>
    {version.name}
    <span class="count">{count}</span>
  </h2>

  <ul>
    {#each Object.entries(games) as game}
      {#if game[1].isInstalled}
        <li>{game[0]}</li>
      {/if}
    {/each}
  </ul>
{/if}
