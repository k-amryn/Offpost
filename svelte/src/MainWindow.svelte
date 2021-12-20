<script lang="typescript">
  import Sidebar from './Sidebar.svelte'
  import Homeview from './Homeview.svelte'
  import Instanceview from './Instanceview.svelte'
  import { activeInstance, ginstances, unsavedChanges } from './stores'

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

  var dialogueActive: boolean = false
  var dialogueType: string
  var instIndex: number = 0
  function showDialogue(type: string, instance: number) {
    dialogueType = type
    dialogueActive = true
    if (instIndex > -1) {
      instIndex = instance
    }
  }

  function deleteInstance(instance: number) {
    $activeInstance = -1
    $unsavedChanges = false

    dialogueActive = false
    var newGinst = $ginstances.slice(0,instance)
    newGinst.push(...$ginstances.slice(instance+1))

    $ginstances = newGinst

    dispatchSocketMessage("d " + instance)

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
    position: relative;
  }
  #mainview {
    overflow: hidden;
  }
  #dialogue {
    position: absolute;
    width: 100%;
    height: 100%;
    background: rgba(0,0,0,0.1);
    display: grid;
    place-items: center;
  }
  #dialogue-inner {
    background: white;
    border-radius: 10px;
    border: var(--main-border-size);
    padding: 20px;
    box-sizing: border-box;
    z-index: 10;
  }
  #delete-content {
    display: grid;
    grid-template-rows: auto 30px;
    width: 300px;
    height: 150px;
  }
  .button-group {
    justify-self: center;
  }
  button {
    width: 70px;
  }
  p {
    margin: 0px 0px 10px 0px;
  }
</style>

<div id="window-inner">
  {#if dialogueActive}
    <div id="dialogue" on:click={() => {dialogueActive = false}}>
      <div id="dialogue-inner" on:click={e => e.stopPropagation()}>
        {#if dialogueType === "delete"}

        <div id="delete-content">
          <div>
            <p>Delete {$ginstances[instIndex].Name}?</p>
            <p>This will remove its post history and platform connections.</p>
          </div>
          <div class="button-group">
            <button on:click={() => deleteInstance(instIndex)}>Yes</button>
            <button on:click={() => {dialogueActive = false}}>No</button>
          </div>
        </div>
        
        {/if}
      </div>
    </div>
  {/if}
  <Sidebar on:alert={(e) => dispatchAlert(e.detail.text)} />

  <div id="mainview">
    {#if $activeInstance === -1}
      <Homeview />
    {:else}
      <Instanceview 
        on:alert={e => dispatchAlert(e.detail.text)}
        on:socketMessage={e => dispatchSocketMessage(e.detail.text)}
        on:dialogue={e => showDialogue(e.detail.type, e.detail.instance)}
      />
    {/if}
  </div>

</div>
