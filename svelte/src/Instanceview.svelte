<script lang="typescript">
  import { createEventDispatcher, onDestroy } from 'svelte'
  import { ginstances, ginstancesOld, activeInstance, unsavedChanges } from './stores'
  import Select from './Select.svelte'

  // Copies config, restores copy when user Cancels changes. Don't copy for
  // newly created instances, because "canceling" should delete new instance
  if ($ginstances[$activeInstance].Status != "new-instance") {
    $ginstancesOld = JSON.parse(JSON.stringify($ginstances))
  }

  // alert about unsaved changes when user clicks "Configure"
  const dispatch = createEventDispatcher()
  function dispatchAlert(msg: string) {
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

  let advanced: boolean = false
  
  $: instance = $ginstances[$activeInstance]
  let platforms: string[] = new Array()
  const unsub = activeInstance.subscribe(() => {
    resetPlatforms()
  })
  onDestroy(unsub)

  function resetPlatforms() {
    if ($activeInstance == -1) return
    platforms = Object.keys($ginstances[$activeInstance].Platforms)
    if (platforms.length == 0) {
      platforms.push("")
    }
    platforms = platforms
  }

  function newPlatform() {
    platforms.push("")
    platforms = platforms
    $unsavedChanges = true
  }

  function removePlatform(i: number) {
    $unsavedChanges = true
    if (i == 0) {
      platforms[i] = ""
      platforms = platforms
    } else {
      platforms.splice(i, 1)
      platforms = platforms
    }
  }

  function addFolder() {
    $unsavedChanges = true
    instance.ImgFolders = [...instance.ImgFolders, ""]
  }

  function removeFolder(i: number) {
    $unsavedChanges = true
    instance.ImgFolders.splice(i, 1)
    instance.ImgFolders = instance.ImgFolders
  }

  function addVariable(variable: string) {
    $unsavedChanges = true
    let captionInput: HTMLInputElement = document.querySelector('.caption-input')!
    switch (variable) {
      case 'filename':
          captionInput.value += "%{filename}";
        break;
      case 'comment':
        captionInput.value += "%{comment}";
        break;
      case 'date':
        captionInput.value += "%{date}"
        break;
    }
  }

  function configurePlatform(platform: string) {
    if ($unsavedChanges) {
      dispatchAlert("Save changes before configuring a platform.")
    } else {
      dispatchDialogue(platform, $activeInstance)
    }
  }

  function changeCounter(target: string, upOrDown: string) {
    $unsavedChanges = true

    if (instance[target].num > 0 && upOrDown === "down") {
      instance[target].num -= 1
    } else if (upOrDown === "up") {
      instance[target].num += 1
    }
  }

  // used to ensure number inputs are only numbers
  function filterNonDigits(e: any, target: string) {
    e.target.value = e.target.value.replace(/\D/, '')
    if (e.target.value === "") {
      instance[target].num = 0
    } else {
      instance[target].num = parseInt(e.target.value, 10)
    }
  }

  function saveInstanceSettings() {
    Object.keys(instance.Platforms).forEach(e => {
      if (!platforms.includes(e)) {
        delete instance.Platforms[e]
      }
    })
    platforms.forEach((e, i) => {
      if (instance.Platforms[e] == undefined && e != "") {
        instance.Platforms[e] = "no-config"
      }
      if (e == "" && platforms.length != 1) {
        removePlatform(i)
      }
    })
    resetPlatforms()

    $unsavedChanges = false
    $ginstancesOld = JSON.parse(JSON.stringify($ginstances))
    dispatchSocketMessage("s " + JSON.stringify($ginstances))
  }

  function cancelInstanceSettings() {
    $unsavedChanges = false
    if ($ginstances[$activeInstance].Status === "new-instance") {
      $activeInstance = $ginstances.length - 2
    } 
    $ginstances = JSON.parse(JSON.stringify($ginstancesOld))
    resetPlatforms()
  }

  function confirmDelete() {
    dispatchDialogue("delete", $activeInstance)
  }
</script>


<style>
  #container {
    box-sizing: border-box;
    height: 100%;
    overflow: auto;
    display: grid;
    grid-template-rows: 1fr 50px;
  }

  #settings {
    padding: 20px;
    overflow: auto;
  }

  #heading-bar {
    border: var(--main-border-size);
    border-radius: 10px;
    display: grid;
    padding: 0px 15px 0px 15px;
    align-items: center;
    grid-template-columns: 1fr;
    font-size: 1em;
  }

  .subtext {
    font-size: var(--small-font);
    text-decoration: underline;
    cursor: pointer;
  }

  .setting-section {
    margin-top: 20px;
    display: grid;
    grid-template-columns: 115px 1fr;
  }

  .setting-label {
    margin-top: 7px;
  }

  #folder-rows {
    width: 100%;
    display: grid;
    grid-template-rows: repeat(auto-fill, 1fr);
    grid-gap: 10px;
  }

  .folder-row {
    display: grid;
    grid-template-columns: 1fr auto auto;
    align-items: center;
    justify-items: start;
    margin-bottom: 0px;
    gap: 20px;
  }

  .svg-holder {
    height: 100%;
    display: grid;
    place-items: center;
    width: 18px;
  }

  .svg-holder:not(.status-indicator) {
    cursor: pointer;
  }

  .folder-input {
    width: 100%;
  }

  #platform-rows {
    width: 100%;
    display: grid;
    grid-template-rows: repeat(auto-fill, 1fr);
    grid-gap: 10px;
  }

  .platform-row {
    display: flex;
    grid-gap: 10px;
  }

  .caption-input {
    width: 100%;
    resize: vertical;
  }

  .caption-subtext {
    font-size: var(--small-font);
    display: block;
  }

  .caption-subtext span {
    cursor: pointer;
    text-decoration: underline;
  }

  .counter {
    display: grid;
    grid-template-columns: 36px 46px 36px;
    align-items: center;
    width: min-content;
    padding: none;
  }

  .counter input {
    text-align: center;
    border-radius: 0px !important;
  }

  .post-delay, .queue-delay {
    display: grid;
    grid-template-columns: min-content min-content;
    grid-gap: 10px;
  }

  .minus, .plus {
    box-sizing: border-box;
    height: 36px;
    display: grid;
    place-items: center;
    cursor: pointer;
    border-top: 2px solid black;
    border-bottom: 2px solid black;
  }

  .minus {
    border-left: 2px solid black;
    border-top-left-radius: 5px;
    border-bottom-left-radius: 5px;
  }

  .plus {
    border-right: 2px solid black;
    border-bottom-right-radius: 5px;
    border-top-right-radius: 5px;
  }

  .advanced {
    margin-top: 20px;
    margin-bottom: 20px;
    display: grid;
    grid-template-columns: min-content min-content;
    align-items: center;
    grid-gap: 6px;
    width: min-content;
    font-size: var(--small-font);
    height: 20px;
    cursor: pointer;
  }

  .delete-instance {
    font-size: var(--small-font);
    color: var(--red);
    cursor: pointer;
    margin-top: 20px;
    display: inline-block;
  }
  
  #status-bar {
    display: grid;
    grid-template-columns: 1fr auto;
    align-items: center;
    padding: 0px 20px 0px 20px;
    box-sizing: border-box;
  }

  #status-bar span {
    font-size: var(--small-font)
  }

  #status-bar button {
    width: 80px;
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

  .status-indicator:not(.not-configured) svg {
    fill: var(--green);
  }

  .status-indicator.not-configured svg {
    fill: var(--orange);
  }
</style>

<div id="container">
  <div id="settings">
    <div id="heading-bar">
      <p id="heading-text">{ instance.Name }</p>
    </div>
    <div class="setting-section">
      <div class="setting-label">Folders:</div>
      <div class="setting-content folders">
        <div id="folder-rows">
          {#each instance.ImgFolders as folder, i}
          <div class="folder-row">
              <input class="folder-input" bind:value={folder} on:input={() => $unsavedChanges = true}>
              <div class="svg-holder">
                <svg width="18px" style="margin-bottom: 2px;" viewBox="0 0 25 21" xmlns="http://www.w3.org/2000/svg" fill-rule="evenodd" clip-rule="evenodd" stroke-linecap="round" stroke-linejoin="round" stroke-miterlimit="1.5">
                  <path d="M2412 271v31h47v-26h-25l-4-5h-18z" fill="none" stroke="#000" stroke-width="5.8" transform="matrix(.45762 0 0 .56483 -1102 -151)"/>
                </svg>
              </div>
              <div class="svg-holder" on:click={() => removeFolder(i)}>
                <svg width="12px" version="1.1" viewBox="0 0 5.4152 5.4152" xmlns="http://www.w3.org/2000/svg">
                  <g transform="translate(-322.81 -103.89)" fill="none" stroke="#de3e39" stroke-linecap="round">
                  <path d="m323.31 104.39 4.4152 4.4152"/>
                  <path d="m327.73 104.39-4.4152 4.4152"/>
                  </g>
                </svg>
              </div>
            </div>
          {/each}
        </div>
        <span on:click={() => addFolder()} class="subtext">Add new folder</span>

      </div>
    </div>

    <div class="setting-section">
      <div class="setting-label">Platforms:</div>
      <div class="setting-content">
        <div id="platform-rows">
          {#each platforms as platform, i}
            <div class="platform-row">
              <Select bind:value={platform} on:change={() => $unsavedChanges=true} width={"155px"} 
                values={{"": "Select Platform","twitter": "Twitter","facebook": "Facebook","tumblr": "Tumblr"}} />
              {#if (platform != "" && instance.Platforms[platform] != undefined)}
                <button on:click={() => configurePlatform(platform)}>Configure</button>
                <div class:not-configured="{instance.Platforms[platform] == "no-config"}" class="svg-holder status-indicator">
                  <svg width="12px" viewBox="0 0 15 15" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" xml:space="preserve">
                    <g transform="matrix(1,0,0,1,-1223.59,-1560.19)">
                    <g transform="matrix(1,0,0,1,0,1101.51)">
                    <g transform="matrix(0.521739,0,0,0.521739,591.235,-354.719)">
                    <circle cx="1225.71" cy="1572.7" r="13.69" />
                    </g>
                    </g>
                    </g>
                  </svg>
                </div>
                {#if (instance.Platforms[platform] != undefined && instance.Platforms[platform] != "no-config")}
                  <div on:click={() => window.open(instance.Platforms[platform])} class="svg-holder">
                    <svg width="14px" style="fill-rule:evenodd;clip-rule:evenodd;stroke-linecap:round;stroke-linejoin:round;stroke-miterlimit:1.5;" viewBox="0 0 15 15" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" xml:space="preserve">
                      <g transform="matrix(1,0,0,1,-1254.39,-1559.84)">
                      <g transform="matrix(1,0,0,1,0,1101.51)">
                      <g transform="matrix(1,0,0,1,22.3939,-5.68736)">
                      <g transform="matrix(6,0,0,6,-6195,-9118.05)">
                      <path d="M1239,1597.34L1238.34,1597.34C1238.15,1597.34 1238,1597.49 1238,1597.68L1238,1599.01C1238,1599.19 1238.15,1599.34 1238.33,1599.34L1239.67,1599.34C1239.85,1599.34 1240,1599.19 1240,1599.01C1240,1598.72 1240,1598.34 1240,1598.34" style="fill:none;stroke:black;stroke-width:0.33px;"/>
                      </g>
                      <g transform="matrix(1,0,0,1,0,-1126.34)">
                      <path d="M1239,1598.34L1245.98,1591.36L1242,1591.36L1245.98,1591.36L1245.98,1595.34" style="fill:none;stroke:black;stroke-width:2px;"/>
                      </g>
                      </g>
                      </g>
                      </g>
                    </svg>
                  </div>
                {/if}
              {/if}
              {#if (i == 0 && (platform != "" && (instance.Platforms[platform] == "no-config" || instance.Platforms[platform] == undefined))) ||
                (i > 0 && (instance.Platforms[platform] == "no-config" || instance.Platforms[platform] == undefined))}
                <div class="svg-holder" on:click={() => removePlatform(i)}>
                  <svg width="12px" version="1.1" viewBox="0 0 5.4152 5.4152" xmlns="http://www.w3.org/2000/svg">
                    <g transform="translate(-322.81 -103.89)" fill="none" stroke="#de3e39" stroke-linecap="round">
                    <path d="m323.31 104.39 4.4152 4.4152"/>
                    <path d="m327.73 104.39-4.4152 4.4152"/>
                    </g>
                  </svg>
                </div>
              {/if}
            </div>
          {/each}
        </div>
        <span on:click={() => newPlatform()} class="subtext">Add new platform</span>
      </div>
    </div>

    <div class="setting-section">
      <div class="setting-label">Caption:</div>
      <div class="setting-content">
        <textarea class="caption-input" rows="2" on:input={() => $unsavedChanges = true} bind:value={instance.Caption}></textarea>
        <span class="caption-subtext">Add variable:
          <span on:click={() => addVariable('filename')}>Filename</span>, 
          <span on:click={() => addVariable('comment')}>Comment</span>,
          <span on:click={() => addVariable('date')}>Date</span>
        </span>
      </div>
    </div>

    <div class="setting-section">
      <div class="setting-label">Post Delay:</div>
      <div class="setting-content post-delay">
        <div class="counter">
          <div on:click={() => changeCounter("PostDelay", "down")} class="minus">
            <svg width="15px" version="1.1" viewBox="0 0 11.863 11.863" xmlns="http://www.w3.org/2000/svg">
              <g transform="translate(-173.99 -176.63)" fill="none" stroke="#000" stroke-linecap="round" stroke-width="1.5218px">
                <path d="m185.09 182.56h-10.341"/>
              </g>
            </svg>
          </div>
          <input on:input={e => {filterNonDigits(e, "PostDelay"); $unsavedChanges = true} } value={ instance.PostDelay.num }>
          <div on:click={() => changeCounter("PostDelay", "up") } class="plus">
            <svg width="15px" version="1.1" viewBox="0 0 11.863 11.863" xmlns="http://www.w3.org/2000/svg">
              <g transform="translate(-173.99 -176.63)" fill="none" stroke="#000" stroke-linecap="round" stroke-width="1.5218px">
                <path d="m179.92 177.39v10.341"/>
                <path d="m185.09 182.56h-10.341"/>
              </g>
            </svg>
          </div>
        </div>
        <Select on:change={() => $unsavedChanges = true} bind:value={ instance.PostDelay.unit } width={"100px"}
          values={{"minutes": "minutes", "hours": "hours", "days": "days"}}
        />
      </div>
    </div>

    <div on:click={() => {advanced = !advanced}} class="advanced">
      <span>Advanced</span>
      {#if !advanced}
        <svg width="9px" viewBox="0 0 111 67" xmlns="http://www.w3.org/2000/svg" fill-rule="evenodd" clip-rule="evenodd" stroke-linecap="round" stroke-linejoin="round" stroke-miterlimit="1.5">
          <path d="M99 11L55 55 11 11" fill="none" stroke="#000" stroke-width="22.4"/>
        </svg>
      {:else}
        <svg width="9px" clip-rule="evenodd" fill-rule="evenodd" stroke-linecap="round" stroke-linejoin="round" stroke-miterlimit="1.5" version="1.1" viewBox="0 0 111 67" xmlns="http://www.w3.org/2000/svg">
          <path d="m11 55 44-44 44 44" fill="none" stroke="#000" stroke-width="22.4"/>
        </svg>     
      {/if}
    </div>

    {#if advanced}
      <div class="setting-section">
        <div class="setting-label">Queue Delay:</div>
        <div class="setting-content queue-delay">
          <div class="counter">
            <div on:click={() => changeCounter("QueueDelay", "down")} class="minus">
              <svg width="15px" version="1.1" viewBox="0 0 11.863 11.863" xmlns="http://www.w3.org/2000/svg">
                <g transform="translate(-173.99 -176.63)" fill="none" stroke="#000" stroke-linecap="round" stroke-width="1.5218px">
                  <path d="m185.09 182.56h-10.341"/>
                </g>
              </svg>
            </div>
            <input on:input={e => {filterNonDigits(e, "QueueDelay"); $unsavedChanges = true} } style="border-radius: 0px;" value={ instance.QueueDelay.num }>
            <div on:click={() => changeCounter("QueueDelay", "up")} class="plus">
              <svg width="15px" version="1.1" viewBox="0 0 11.863 11.863" xmlns="http://www.w3.org/2000/svg">
                <g transform="translate(-173.99 -176.63)" fill="none" stroke="#000" stroke-linecap="round" stroke-width="1.5218px">
                  <path d="m179.92 177.39v10.341"/>
                  <path d="m185.09 182.56h-10.341"/>
                </g>
              </svg>
            </div>
          </div>
          <Select on:change={() => $unsavedChanges = true} bind:value={ instance.QueueDelay.unit } width="100px"
            values={{"seconds": "seconds", "minutes": "minutes"}}/>
        </div>
      </div>

      <div class="setting-section">
        <div class="setting-label">Startup Delay:</div>
        <div class="setting-content">
          <!-- On startup, attempt posting at the NextPostTime. If the NextPostTime has passed, 
          then use this option. -->
          <Select on:change={() => $unsavedChanges = true} bind:value={ instance.StartupPostDelay } width={"100px"}
            values={{ "random": "Random", "full": "Full", "none": "None" }} />
        </div>
      </div>

      <div class="delete-instance">
        <span on:click={() => confirmDelete()}>Delete instance</span>
      </div>
    {/if}
  </div>
  <div id="status-bar">
    {#if $unsavedChanges}
      <span class="needs-configuring">Configuring...</span>
    {:else if instance.Status === "needs-configuring"}
      <span class="needs-configuring">Needs configuring</span>
    {:else if instance.ItemsInQueue > 0}
      <span class="queued">{instance.ItemsInQueue} items in queue<br>
        Next post {instance.NextPostTime}<br></span>
    {:else if instance.ItemsInQueue === 0}
      <span class="waiting">Waiting for new image</span>
    {/if}
    <div class="button-group">
      {#if !$unsavedChanges}
      <button>Pause</button>
      {:else}
      <button on:click={saveInstanceSettings}>Save</button>
      <button on:click={cancelInstanceSettings}>Cancel</button>
      {/if}
    </div>
  </div>

</div>
