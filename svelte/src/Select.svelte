<script lang="typescript">
  export let value: string = ""
  export let values: object
  $: keys = Object.keys(values)

  let open = false
  export let width: string
  import { createEventDispatcher } from 'svelte'
  const dispatch = createEventDispatcher()
  function dispatchChange(value: string) {
    dispatch('change', {
      value: value
    })
  }

  function handleWindowEvent(e) {
    if (!e.target.classList.contains('select')) {
      open = false
    }
  }

  function selectItem(item: string) {
    dispatchChange(value)
    value = item
    open = false
  }

</script>

<style>
  #container {
    position: relative;
    user-select: none;
    margin-right: 5px;
    -webkit-user-select: none;
  }

  #chosen {
    background: url('data:image/svg+xml;base64,PHN2ZyB2aWV3Qm94PSIwIDAgMTExIDY3IiB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIGZpbGwtcnVsZT0iZXZlbm9kZCIgY2xpcC1ydWxlPSJldmVub2RkIiBzdHJva2UtbGluZWNhcD0icm91bmQiIHN0cm9rZS1saW5lam9pbj0icm91bmQiIHN0cm9rZS1taXRlcmxpbWl0PSIxLjUiPjxwYXRoIGQ9Ik05OSAxMUw1NSA1NSAxMSAxMSIgZmlsbD0ibm9uZSIgc3Ryb2tlPSIjMDAwIiBzdHJva2Utd2lkdGg9IjIyLjQiLz48L3N2Zz4='), var(--filler);
    background-size: 10px;
    background-position: calc(100% - 6px) center;
    background-repeat: no-repeat;
	  cursor: pointer;
    height: 36px;
    box-sizing: border-box;
    padding: 0px 20px 0px 8px;
    border-radius: 5px;
    display: grid;
    align-items: center;
  }

  #chosen.open {
    border: 2px solid black;
    padding: 6px;
    background-position: calc(100% - 4px) center;
  }

  #dropdown {
    position: absolute;
    margin-top: -4px; 
    box-sizing: border-box;
    padding: 3px 3px 6px 3px;
    width: 100%;
    background: var(--filler);
    border-bottom-right-radius: 5px;
    border-bottom-left-radius: 5px;
    border-left: 2px solid black;
    border-right: 2px solid black;
    border-bottom: 2px solid black;
    z-index: 2;

  }

  .value {
    height: 30px;
    box-sizing: border-box;
    padding: 0px 3px 0px 3px;
    display: grid;
    align-items: center;
    cursor: default;
  }

  .value:hover {
    background: rgba(0,0,0,0.1);
    border-radius: 5px;
  }
</style>

<svelte:window on:mousedown={handleWindowEvent} />

<div id="container" class="select" style="width: {width}">
  <div id="chosen" class="select" class:open on:mousedown={() => {open = !open}}>{values[value]}</div>
  {#if open}
  <div class:visible={open} id="dropdown" class="select">
    {#each keys as key}
      {#if key != value}
        <div class="select value" on:mouseup={() => selectItem(key)}>{values[key]}</div>
      {/if}
    {/each}
  </div>
  {/if}
</div>