<script lang="ts">
	import MainWindow from './MainWindow.svelte';
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
