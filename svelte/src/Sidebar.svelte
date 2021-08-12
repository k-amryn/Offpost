<script lang="ts">
	let list = [
		"1",
		"2",
		"3",
		"4",
		"5",
		"6",
		"7"
	]
</script>

<style>
	#container {
		border-right: 5px solid black;

	}
	.item {
		display: grid;
		place-items: center;
	}
</style>

<div id="container">
	{#each list as item (item)}
		<div class="item">{item}</div>
	{/each}
</div>
