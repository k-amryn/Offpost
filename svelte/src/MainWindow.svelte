<script lang="typescript">
  import Sidebar from './Sidebar.svelte'
  import Homeview from './Homeview.svelte'
  import Instanceview from './Instanceview.svelte'
  import { activeInstance} from './stores'

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

  function dispatchDialogue(type: string, instance: number) {
    dispatch('dialogue', {
      type: type,
      instance: instance
    })
  }
</script>

<style>
  #window-inner {
    background: white;
    border-radius: 20px;
    width: 100%;
    min-height: 700px;
    display: grid;
    grid-template-columns: auto 1fr;
    gap: 15px;
    overflow: hidden;
    box-sizing: border-box;
    position: relative;
    padding: 15px;
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
        on:dialogue={e => dispatchDialogue(e.detail.type, e.detail.instance)}
      />
    {/if}
  </div>

</div>
