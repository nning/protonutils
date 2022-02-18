<script>
  import Tabs from './tabs/Tabs.svelte';

  import ListTab from './tabs/list/ListTab.svelte';
  import TestTab from './tabs/test/TestTab.svelte';
</script>

<style lang="scss">
  @import './main.scss';

  :global(body) {
    color: $foreground-lightest;
    height: 100vh;
    background: radial-gradient(circle at top left, rgba($background-light, 1) 0%, rgba($background-dark, 1) 100%);
  }
</style>

<Tabs
  tabs={[
    {title: 'List', component: ListTab},
    {title: 'Test', component: TestTab}
  ]}
/>
