<script>
  export let version = {};

  $: games = version?.games && Object.values(version?.games)

  function hasInstalledGames() {
    try {
      return !!games.find(game => game.isInstalled);
    } catch(e) {
      return false;
    }
  }
</script>

{#if hasInstalledGames()}
  <h2>{version.name}, {games.length}</h2>
  <ul>
    {#each games as game}
      {#if game.isInstalled}
        <li>{game.name}</li>
      {/if}
    {/each}
  </ul>
{/if}
