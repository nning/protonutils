<script>
	import { onMount } from "svelte";

	let output = '';

	onMount(() => {
		window.backend.list().then(result => {
      // result = JSON.parse(result)
			output = result;
		});
	});
</script>

<div class="App">
	<h1>Hello World!</h1>
	<h2>{output}</h2>
</div>
