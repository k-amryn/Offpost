<script lang="typescript">
  import { ginstances } from './stores'
  export let instance: number
  $: connected = $ginstances[instance].Platforms["twitter"] != "no-config"

  import { createEventDispatcher } from 'svelte'
  const dispatch = createEventDispatcher()
  function dispatchSocketMessage(msg: string) {
    dispatch('socketMessage', {
      text: msg
    })
  }

  let confirming: boolean = false
</script>

<style>
  img {
    pointer-events: none;
    user-select: none;
    -webkit-user-drag: none;
  }

  svg {
    stroke: var(--black);
    stroke-width: 160;
  }
    
  svg.connected {
    stroke: var(--green);
    animation: breathe 1s infinite;
    animation-direction: alternate;
  }

  @keyframes breathe {
    from {
      stroke-width: 180;
      -webkit-filter: drop-shadow( 0px 0px 5px var(--green));
      filter: drop-shadow( 0px 0px 5px var(--green));
    }
    to {
      stroke-width: 130;
      -webkit-filter: drop-shadow( 0px 0px 3px var(--green));
      filter: drop-shadow( 0px 0px 3px var(--green));
    }
  }

  #wrapper {
    height: 300px;
    width: 450px;
    display: grid;
    align-items: center;
    position: relative;
  }

  #twitter-header {
    display: grid;
    grid-template-columns: 100px 70px 100px;
    place-content: center;
    height: 200px;
    grid-gap: 8px;
  }

  #arrow {
    display: grid;
    place-items: center;
  }

  #content {
    box-sizing: border-box;
    height:60px;
    width: 100%;
    position: absolute;
    bottom: 10px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: end;
  }

  .button-line {
    width: 100%;
    display: flex;
    justify-content: center;
    gap: 10px;
  }

  .button-line.connect {
    justify-self: start;
  }

  .confirmation button {
    width: 80px;
  }
  p {
    margin: none;
  }

  .disconnect button {
    border: none;
    color: var(--red);
  }
</style>

<div id="wrapper">
  <div id="twitter-header">
    <img width="100px" src="logo.svg" alt="">
    <div id="arrow">
      <svg class:connected viewBox="0 0 1086 473" xmlns="http://www.w3.org/2000/svg" xml:space="preserve" fill-rule="evenodd" clip-rule="evenodd" stroke-linecap="round" stroke-linejoin="round" stroke-miterlimit="1.5">
        <path d="m1955.4 2209.9 413.4 426.4-413.4-426.4 413.4-426.5-413.4 426.5h2128L3670 1783.4l413.4 426.5-413.4 426.4 413.4-426.4h-2128Z" fill="none" transform="matrix(.47099 0 0 .45657 -879.5 -772.7)"/>
      </svg>
    </div>
    <img style="margin-top: 5px" width="100px" src="twitter.svg" alt="">
  </div>


  <div id="content">
    {#if !connected}
      <div class="button-line connect">
        <button on:click={() => dispatchSocketMessage("twitter " + instance)}>Connect</button>
      </div>
    {:else}
      {#if !confirming}
        <div class="button-line disconnect">
          <button on:click={() => confirming = true}>Disconnect</button>
        </div>
      {:else}
        <p>Are you sure? You can reconnect later.</p>
        <div class="button-line confirmation">
          <button on:click={() => dispatchSocketMessage("twitter,r " + instance)}>Yes</button>
          <button on:click={() => confirming = false}>No</button>
        </div>
      {/if}
    {/if}
  </div>
</div>