<script>
  import Tabs from './tabs/Tabs.svelte';

  import ListTab from './tabs/list/ListTab.svelte';
  import TestTab from './tabs/test/TestTab.svelte';
</script>

<div class="App">
  <Tabs
    tabs={[
      {title: 'List', component: ListTab},
      {title: 'Test', component: TestTab}
    ]}
    />
</div>
