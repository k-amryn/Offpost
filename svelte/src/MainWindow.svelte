<script lang="typescript">
  import Sidebar from './Sidebar.svelte'
  import Homeview from './Homeview.svelte'
  import { instances, activeInstance } from './stores'
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
  <Sidebar />

  <div id="mainview">
    {#if $activeInstance === -1}
      <Homeview />
    {/if}
  </div>

</div>
