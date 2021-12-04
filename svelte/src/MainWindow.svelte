<script lang="typescript">
  import Sidebar from './Sidebar.svelte'
  import Homeview from './Homeview.svelte'
  import Instanceview from './Instanceview.svelte'
  import { activeInstance } from './stores'

  // relay alert messages from sub-components to App component
  import { createEventDispatcher } from 'svelte'
  const dispatch = createEventDispatcher()
  function dispatchAlert(msg) {
    dispatch('alert', {
      text: msg
    })
  }

  function dispatchSocketMessage(msg: string) {
    dispatch('socketMessage', {
      text: msg
    })
  }
</script>

<style>
  #window-inner {
    border: 5px solid black;
    border-radius: 10px;
    width: 100%;
    height: 700px;
    display: grid;
    grid-template-columns: 100px 1fr;
    overflow: hidden;
    box-sizing: border-box;
  }
  #mainview {
    overflow: hidden;
  }
</style>

<div id="window-inner">
  <Sidebar on:alert={(e) => dispatchAlert(e.detail.text)} />

  <div id="mainview">
    {#if $activeInstance === -1}
      <Homeview />
    {:else}
      <Instanceview 
        on:alert={e => dispatchAlert(e.detail.text)}
        on:socketMessage={e => dispatchSocketMessage(e.detail.text)}
      />
    {/if}
  </div>

</div>
