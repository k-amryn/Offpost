<script lang="typescript">
	import { ginstances, activeInstance } from './stores'

	function setActive(instanceIndex: number) {
		$activeInstance = instanceIndex
	}
</script>

<style>
	#container {
		box-sizing: border-box;
		height: 100%;
		overflow: auto;
	}
	

	#heading-bar {
		background: var(--filler);
		border-radius: 15px;
		height: 90px;
		padding: 0px 20px 0px 20px;
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 1em;
		position: relative;
		transition: background 0.2s;
	}

	#heading-text {
		height: 100%;
		display: flex; 
		place-items: center;
		place-content: center;
	}

	#heading-menu {
		position: absolute;
		right: 20px;
		width: 20px;
		cursor: pointer;
	}

	#content {
		margin-top: 10px;
		display: grid;
		grid-template-columns: 0.5fr 0.5fr;
		width: 100%;
		gap: 10px;
	}

	.instance {
		background: var(--filler);
		border-radius: 10px;
		padding: 10px;
		display: grid;
		grid-template-columns: auto 1fr;
		gap: 10px;
		cursor: pointer;
		user-select: none;
	}

	.instance-img {
		width: 100px;
		height: 100px;
		border-radius: 7px;
		overflow: hidden;
	}

	.instance-img img {
		width: 100%;
		height: 100%;
		object-fit: cover;
	}

	.instance-text {
		width: 192px;
	}

	.instance-status span {
		font-size: var(--small-font);
	}

	.instance-name {
		text-overflow: ellipsis;
		overflow-x: hidden;
		white-space: nowrap;
	}

	span.waiting {
		color: var(--blue);
	}

	span.queued {
		color: var(--green);
	}

	span.needs-configuring {
		color: var(--orange);
	}
</style>

<div id="container">
	<div id="heading-bar">
		<div id="heading-text">
			<img src="images/cottage.svg" alt="">
		</div>
		<svg id="heading-menu" version="1.1" viewBox="0 0 7.8668 7.5988" xmlns="http://www.w3.org/2000/svg">
			<g transform="translate(-319.08 -86.143)">
			<g transform="translate(149.08 67.09)">
			<g fill="none" stroke="#000" stroke-linecap="round" stroke-width="1.0104">
			<path d="m170.5 19.558h6.8564"/>
			<path d="m170.5 22.852h6.8564"/>
			<path d="m170.5 26.146h6.8564"/>
			</g>
			</g>
			</g>
		</svg>
	</div>
	<div id="content">
		{ #each $ginstances as instance, i (instance.Name) }
			<div class="instance" on:click="{() => setActive(i)}">

				<div class="instance-img">
					<img src="{instance.Image}" alt="{instance.Name}">
				</div>

				<div class="instance-text">
					<div class="instance-name">{instance.Name}</div>
					<div class="instance-status">
						{#if instance.Status === "needs-configuring"}
							<span class="needs-configuring">Needs configuring</span>
						{:else if instance.ItemsInQueue > 0}
							<span class="queued">{instance.ItemsInQueue} items in queue<br>
								Next post {instance.NextPostTime}<br></span>
						{:else if instance.ItemsInQueue === 0}
							<span class="waiting">Waiting for new image</span>
						{/if}
					</div>
				</div>

			</div>
		{ /each }
	</div>
</div>
