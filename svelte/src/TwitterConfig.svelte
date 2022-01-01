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

  #connected {
    color: green;
  }

  #twitter-header {
    margin-top: 20px;
    margin-bottom: 20px;
    display: grid;
    grid-template-columns: 100px 70px 100px;
    place-content: center;
    grid-gap: 8px;
  }

  #arrow {
    display: grid;
    place-items: center;
  }

  #button-line {
    width: 100%;
    text-align: center;
  }
</style>

<div id="wrapper">
  <div id="twitter-header">
    <img width="100px" src="logo.svg" alt="">
    <div id="arrow">
      <svg viewBox="0 0 1086 473" xmlns="http://www.w3.org/2000/svg" xml:space="preserve" fill-rule="evenodd" clip-rule="evenodd" stroke-linecap="round" stroke-linejoin="round" stroke-miterlimit="1.5">
        <path d="m1955.4 2209.9 413.4 426.4-413.4-426.4 413.4-426.5-413.4 426.5h2128L3670 1783.4l413.4 426.5-413.4 426.4 413.4-426.4h-2128Z" fill="none" stroke="#000" stroke-width="158.9" transform="matrix(.47099 0 0 .45657 -879.5 -772.7)"/>
      </svg>
    </div>
    <img style="margin-top: 5px" width="100px" src="twitter.svg" alt="">
  </div>

  <div id="button-line">
    <button on:click={() => dispatchSocketMessage("twitter " + instance)}>Connect</button>
  </div>

  {#if ($ginstances[instance].Platforms["twitter"] != "no" && $ginstances[instance].Platforms["twitter"] != "no-config")}
    <p>Status: <span id="connected">Connected</span></p>
  {/if}
</div>