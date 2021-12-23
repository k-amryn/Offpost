<script lang="typescript">
  import { ginstances } from './stores'
  export let instance: number

  import { createEventDispatcher } from 'svelte'
  const dispatch = createEventDispatcher()
  function dispatchSocketMessage(msg: string) {
    dispatch('socketMessage', {
      text: msg
    })
  }
</script>

<style>
  #wrapper {
    height: 500px;
    width: 450px;

  }
</style>

<div id="wrapper">
  <p>Here we're configuring {$ginstances[instance].Name} for twitter.</p>
  <button on:click={() => dispatchSocketMessage("twitter " + instance)}>Connect</button>
</div>