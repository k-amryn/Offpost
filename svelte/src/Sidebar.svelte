<script lang="typescript">
  import { ginstances, activeInstance, unsavedChanges } from './stores'

  // alert about unsaved changes when user tries to switch instance
  import { createEventDispatcher } from 'svelte'
  const dispatch = createEventDispatcher()
  function dispatchAlert(msg) {
    dispatch('alert', {
      text: msg
    })
  }

  function setActive(instanceIndex: number) {
    if ($unsavedChanges) {
      dispatchAlert("Save changes before switching instances.")
    } else {
      $activeInstance = instanceIndex;
    }
  }

  function newInstance() {
    // new instances have names with increasing numbers
    var name: string
    var baseTaken: boolean = false
    var unusedNum: number = 2
    for (let i = 0; i < $ginstances.length; i++) {
      if (baseTaken) {
        if ($ginstances[i].Name == "New Instance (" + unusedNum.toString() + ")") {
          unusedNum++
          i = 0
        }
      } else if ($ginstances[i].Name == "New Instance") {
        baseTaken = true
        i = 0
      }
    }

    if (baseTaken) {
      name = "New Instance (" + unusedNum.toString() + ")"
    } else {
      name = "New Instance"
    }

    if (!$unsavedChanges) {
      $ginstances = [...$ginstances, {
        Name: name,
        ImgFolders: [""],
        QueueDelay: {num: 30, unit: "seconds"},
        PostDelay: {num: 5, unit: "hours"},
        StartupPostDelay: "random",
        Platforms: {},
        Caption: "",
        ItemsInQueue: 0,
        NextPostTime: "",
        Status: "new-instance",
        Image: "./new_instance.svg"
      }]
    }
    setActive($ginstances.length - 1)
    $unsavedChanges = true
  }
</script>

<style>
  #sidebar-items-and-space {
    display: grid;
    grid-template-rows: min-content 1fr;
    place-items: center;
    justify-items: center;
    overflow: auto;
    scrollbar-width: none;
    box-sizing: border-box;
    height: 100%;
  }

  #sidebar-items-and-space::-webkit-scrollbar {
    display: none;
  }

  #afterspace {
    border-right: var(--main-border-size);
    border-top: none;
    height: 100%;
    width: 100%;
    box-sizing: border-box;
  }

  #sidebar {
    width: 100%;
    display: grid;
    place-items: center;
  }

  .sidebar-item {
    width: 100%;
    height: 90px;
    display: grid;
    place-items: center;
    box-sizing: border-box;
    overflow: hidden;
  }

  .sidebar-item-inner {
    width: 70px;
    height: 70px;
    border: 2px solid black;
    border-radius: 7px;
    overflow: hidden;
    display: grid;
    place-items: center;
  }

  #home .sidebar-item-inner {
    border: none;
  }

  #new-instance .sidebar-item-inner {
    border: none;
  }

  #home svg, .sidebar-item img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .sidebar-item.active {
    margin-left: -5px;
  }

  .sidebar-item:not(.active) {
    cursor: pointer;
    border-right: var(--main-border-size);
  }

  .sidebar-borderer {
    height: 10px;
    width: 100%;
    box-sizing: border-box;
    border-right: var(--main-border-size);
  }

  .before-active {
    border-bottom-right-radius: 10px;
    border-bottom: var(--main-border-size);
  }

  .after-active {
    border-top-right-radius: 10px;
    border-top: var(--main-border-size);
  }
</style>

<div id="sidebar-items-and-space">
  <div id="sidebar">
    <div class="none"></div>

    <div id="home" class="sidebar-item"
      class:active="{ $activeInstance === -1 }"
      on:click={() => setActive(-1)}>
      <div class="sidebar-item-inner">
        <svg version="1.1" viewBox="0 0 15.524 15.749" xmlns="http://www.w3.org/2000/svg">
          <g transform="translate(-183.18 -96.312)">
          <g transform="translate(10.716 13.759)">
          <g transform="translate(148.38 35.711)">
          <rect transform="scale(1,-1)" x="24.39" y="-54.292" width="6.9959" height="6.9959" rx=".90915"/>
          <rect transform="scale(1,-1)" x="32.312" y="-62.137" width="6.9959" height="6.9959" rx=".90915"/>
          <rect transform="scale(1,-1)" x="24.39" y="-62.137" width="6.9959" height="6.9959" rx=".90915"/>
          <rect transform="scale(1,-1)" x="32.312" y="-54.292" width="6.9959" height="6.9959" rx=".90915"/>
          </g>
          </g>
          </g>
        </svg>
      </div>
    </div>

    { #each $ginstances as instance, index (instance.Name)}
      <div class="sidebar-borderer"
           class:before-active="{ index === $activeInstance }"
           class:after-active="{ index === $activeInstance + 1 }"></div>
      <div class="sidebar-item"
        class:active="{ index === $activeInstance }"
        on:click={() => setActive(index)}>
        <div class="sidebar-item-inner">
          <img src="{instance.Image}" alt="{instance.Name} image">
        </div>
      </div>
    { /each }

    <div class="sidebar-borderer" class:after-active="{ $activeInstance === $ginstances.length-1 }"></div>
    <div id="new-instance" class="sidebar-item" on:click="{newInstance}">
      <div class="sidebar-item-inner">
        <svg width="44.835" height="44.835" version="1.1" viewBox="0 0 11.863 11.863" xmlns="http://www.w3.org/2000/svg">
         <g transform="translate(-173.99 -176.63)" fill="none" stroke="#000" stroke-linecap="round" stroke-width="1.5218px">
          <path d="m179.92 177.39v10.341"/>
          <path d="m185.09 182.56h-10.341"/>
         </g>
        </svg>
      </div>
    </div>

  </div>
  <div id="afterspace"></div>
</div>
