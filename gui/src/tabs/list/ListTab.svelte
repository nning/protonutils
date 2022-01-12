<script>
  import { onMount } from 'svelte';

  import Version from './Version.svelte';

  let versions = {};

  onMount(() => {
    window.backend.list().then(result => {
      versions = Object.values(JSON.parse(result))
    });
  });
</script>

<div>
  {#each versions as version (version.id) }
    <Version {version}/>
  {/each}
</div>
