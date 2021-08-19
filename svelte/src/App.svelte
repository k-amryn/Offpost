<script lang="typescript">
	import MainWindow from './MainWindow.svelte';
	import { instances } from './stores'

  fetch("./testinguserdata/offpost.json")
    .then(response => response.json())
    .then(data => { $instances=data })

</script>

<style>
	#the-container {
		width: 800px;
		margin: auto;
		display: grid;
		place-items: center;
	}
</style>

<div id="the-container">
	<div id="brand-heading">
		<img src="./logo.svg" alt="Offpost logo" width="200px">
	</div>
	<MainWindow />
</div>
