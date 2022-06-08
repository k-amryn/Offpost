<script lang="typescript">
  import { createEventDispatcher, onDestroy } from 'svelte'
  import { ginstances, ginstancesOld, activeInstance, unsavedChanges } from './stores'
  import Select from './Select.svelte'

  let advanced: boolean = false
  
  $: instance = $ginstances[$activeInstance]
  let platforms: string[] = new Array()
  const unsub = activeInstance.subscribe(() => {
    resetPlatforms()
  })
  onDestroy(unsub)

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

  function profileLink(platform: string): string {
    switch (platform) {
      case "twitter":
        return "https://twitter.com/" + instance.Platforms["twitter"]
      default:
        break;
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
  }

  #settings {
    overflow: auto;
  }

  #heading-bar {
    background: var(--filler);
    border-radius: 15px;
    height: 90px;
    display: grid;
    padding: 0px 20px 0px 20px;
    align-items: center;
    grid-template-columns: 1fr auto;
    font-size: 1em;
    transition: background 0.2s;
  }

  #heading-bar.unsaved {
    background: rgba(240, 150, 5, 0.25);
  }

  #heading-bar-status {
    display: grid;
    grid-template-columns: 1fr auto;
    align-items: center;
    box-sizing: border-box;
    font-size: var(--small-font);
  }

  .button-group button {
    height: 50px;
    width: 50px;
    border-radius: 10px;
    border: 3px solid black;
    display: grid;
    place-content: center;
    background: transparent;
  }
  
  .subtext {
    font-size: var(--small-font);
    text-decoration: underline;
    cursor: pointer;
  }

  .setting-section {
    margin-top: 20px;
    display: grid;
    border: 3px solid black;
    padding: 10px;
    border-radius: 10px;
    position: relative;
    grid-template-columns: 1fr;
  }

  .setting-label {
    position: absolute;
    background: white;
    padding: 0px 5px 0px 5px;
    font-size: var(--small-font);
    top: -10px;
    left: 10px;
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
    gap: 10px;
  }

  .folder-row input {
    margin-inline-end: 20px;
  }

  .svg-holder {
    height: 100%;
    display: grid;
    place-items: center;
    width: 30px;
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
    grid-gap: 5px;
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
    background: var(--filler);
  }

  .minus {
    border-top-left-radius: 5px;
    border-bottom-left-radius: 5px;
  }

  .plus {
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

  /* .disabled {
    background: rgba(0,0,0,0.1);
    color: rgba(0,0,0,0.5);
    border-color: rgba(0,0,0,0.5);
    cursor: default;
  } */
</style>

<div id="container">
  <div id="settings">
    <div id="heading-bar" class:unsaved={$unsavedChanges}>

      <div id="heading-bar-text">
        <span id="heading-bar-title">{ instance.Name }</span>
        <div id="heading-bar-status">
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
        </div>
      </div>
      
      <div class="button-group">
        {#if !$unsavedChanges}
        <button>
          <svg xmlns="http://www.w3.org/2000/svg" height="40" width="40"><path d="M24.583 31.667Q23.417 31.667 22.604 30.854Q21.792 30.042 21.792 28.875V11.125Q21.792 9.958 22.604 9.146Q23.417 8.333 24.583 8.333H28.875Q30.042 8.333 30.854 9.146Q31.667 9.958 31.667 11.125V28.875Q31.667 30.042 30.854 30.854Q30.042 31.667 28.875 31.667ZM11.125 31.667Q9.958 31.667 9.146 30.854Q8.333 30.042 8.333 28.875V11.125Q8.333 9.958 9.146 9.146Q9.958 8.333 11.125 8.333H15.417Q16.583 8.333 17.396 9.146Q18.208 9.958 18.208 11.125V28.875Q18.208 30.042 17.396 30.854Q16.583 31.667 15.417 31.667Z"/></svg>
        </button>
        {:else}
        <button on:click={saveInstanceSettings}>
          <svg xmlns="http://www.w3.org/2000/svg" height="40" width="40"><path d="M7.792 35Q6.667 35 5.833 34.167Q5 33.333 5 32.208V7.792Q5 6.667 5.833 5.833Q6.667 5 7.792 5H27.25Q27.833 5 28.333 5.208Q28.833 5.417 29.25 5.833L34.167 10.75Q34.583 11.167 34.792 11.667Q35 12.167 35 12.75V32.208Q35 33.333 34.167 34.167Q33.333 35 32.208 35ZM20 29.875Q21.875 29.875 23.208 28.542Q24.542 27.208 24.542 25.292Q24.542 23.417 23.229 22.083Q21.917 20.75 20 20.75Q18.125 20.75 16.792 22.083Q15.458 23.417 15.458 25.292Q15.458 27.208 16.771 28.542Q18.083 29.875 20 29.875ZM11.208 16H23.375Q23.958 16 24.354 15.604Q24.75 15.208 24.75 14.625V11.208Q24.75 10.583 24.354 10.188Q23.958 9.792 23.375 9.792H11.208Q10.583 9.792 10.188 10.188Q9.792 10.583 9.792 11.208V14.625Q9.792 15.208 10.188 15.604Q10.583 16 11.208 16Z"/></svg>
        </button>
        <button on:click={cancelInstanceSettings}>
          <svg xmlns="http://www.w3.org/2000/svg" height="40" width="40"><path d="M20 21.958 11.458 30.5Q11.042 30.917 10.479 30.917Q9.917 30.917 9.5 30.5Q9.083 30.083 9.083 29.521Q9.083 28.958 9.5 28.542L18.042 20L9.5 11.458Q9.083 11.042 9.083 10.479Q9.083 9.917 9.5 9.5Q9.917 9.083 10.479 9.083Q11.042 9.083 11.458 9.5L20 18.042L28.542 9.5Q28.958 9.083 29.521 9.083Q30.083 9.083 30.5 9.5Q30.917 9.917 30.917 10.479Q30.917 11.042 30.5 11.458L21.958 20L30.5 28.542Q30.917 28.958 30.917 29.521Q30.917 30.083 30.5 30.5Q30.083 30.917 29.521 30.917Q28.958 30.917 28.542 30.5Z"/></svg>
        </button>
        {/if}
      </div>

    </div>
    <div class="setting-section">
      <div class="setting-label">Folders</div>
      <div class="setting-content folders">
        <div id="folder-rows">
          {#each instance.ImgFolders as folder, i}
          <div class="folder-row">
              <input class="folder-input" bind:value={folder} on:input={() => $unsavedChanges = true}>
              <div class="svg-holder">
                <svg xmlns="http://www.w3.org/2000/svg" height="24" width="24"><path d="M4 20Q3.175 20 2.588 19.413Q2 18.825 2 18V6Q2 5.175 2.588 4.588Q3.175 4 4 4H9.175Q9.575 4 9.938 4.15Q10.3 4.3 10.575 4.575L12 6H20Q20.825 6 21.413 6.588Q22 7.175 22 8V18Q22 18.825 21.413 19.413Q20.825 20 20 20Z"/></svg>
              </div>
              <div class="svg-holder" on:click={() => removeFolder(i)}>
                <svg xmlns="http://www.w3.org/2000/svg" height="24" width="24"><path d="M7 21Q6.175 21 5.588 20.413Q5 19.825 5 19V6Q4.575 6 4.287 5.713Q4 5.425 4 5Q4 4.575 4.287 4.287Q4.575 4 5 4H9Q9 3.575 9.288 3.287Q9.575 3 10 3H14Q14.425 3 14.713 3.287Q15 3.575 15 4H19Q19.425 4 19.712 4.287Q20 4.575 20 5Q20 5.425 19.712 5.713Q19.425 6 19 6V19Q19 19.825 18.413 20.413Q17.825 21 17 21ZM9 16Q9 16.425 9.288 16.712Q9.575 17 10 17Q10.425 17 10.713 16.712Q11 16.425 11 16V9Q11 8.575 10.713 8.287Q10.425 8 10 8Q9.575 8 9.288 8.287Q9 8.575 9 9ZM13 16Q13 16.425 13.288 16.712Q13.575 17 14 17Q14.425 17 14.713 16.712Q15 16.425 15 16V9Q15 8.575 14.713 8.287Q14.425 8 14 8Q13.575 8 13.288 8.287Q13 8.575 13 9Z"/></svg>
              </div>
            </div>
          {/each}
        </div>
        <span on:click={() => addFolder()} class="subtext">Add new folder</span>

      </div>
    </div>

    <div class="setting-section">
      <div class="setting-label">Platforms</div>
      <div class="setting-content">
        <div id="platform-rows">
          {#each platforms as platform, i}
            <div class="platform-row">
              <Select bind:value={platform} on:change={() => $unsavedChanges=true} width={"155px"} 
                values={{"": "Select Platform","twitter": "Twitter","facebook": "Facebook","tumblr": "Tumblr"}} />
              {#if platform != ""}
                <div on:click={() => configurePlatform(platform)} class="svg-holder configure-button">
                  <svg xmlns="http://www.w3.org/2000/svg" height="24" width="24"><path d="M13.875 22H10.125Q9.75 22 9.475 21.75Q9.2 21.5 9.15 21.125L8.85 18.8Q8.525 18.675 8.238 18.5Q7.95 18.325 7.675 18.125L5.5 19.025Q5.15 19.15 4.8 19.05Q4.45 18.95 4.25 18.625L2.4 15.4Q2.2 15.075 2.275 14.7Q2.35 14.325 2.65 14.1L4.525 12.675Q4.5 12.5 4.5 12.337Q4.5 12.175 4.5 12Q4.5 11.825 4.5 11.662Q4.5 11.5 4.525 11.325L2.65 9.9Q2.35 9.675 2.275 9.3Q2.2 8.925 2.4 8.6L4.25 5.375Q4.425 5.025 4.787 4.937Q5.15 4.85 5.5 4.975L7.675 5.875Q7.95 5.675 8.25 5.5Q8.55 5.325 8.85 5.2L9.15 2.875Q9.2 2.5 9.475 2.25Q9.75 2 10.125 2H13.875Q14.25 2 14.525 2.25Q14.8 2.5 14.85 2.875L15.15 5.2Q15.475 5.325 15.763 5.5Q16.05 5.675 16.325 5.875L18.5 4.975Q18.85 4.85 19.2 4.95Q19.55 5.05 19.75 5.375L21.6 8.6Q21.8 8.925 21.725 9.3Q21.65 9.675 21.35 9.9L19.475 11.325Q19.5 11.5 19.5 11.662Q19.5 11.825 19.5 12Q19.5 12.175 19.5 12.337Q19.5 12.5 19.45 12.675L21.325 14.1Q21.625 14.325 21.7 14.7Q21.775 15.075 21.575 15.4L19.725 18.6Q19.525 18.925 19.163 19.038Q18.8 19.15 18.45 19.025L16.325 18.125Q16.05 18.325 15.75 18.5Q15.45 18.675 15.15 18.8L14.85 21.125Q14.8 21.5 14.525 21.75Q14.25 22 13.875 22ZM12.05 15.5Q13.5 15.5 14.525 14.475Q15.55 13.45 15.55 12Q15.55 10.55 14.525 9.525Q13.5 8.5 12.05 8.5Q10.575 8.5 9.562 9.525Q8.55 10.55 8.55 12Q8.55 13.45 9.562 14.475Q10.575 15.5 12.05 15.5Z"/></svg>
                </div>
                <div class:not-configured="{(instance.Platforms[platform] == "no-config" || instance.Platforms[platform] == undefined)}" class="svg-holder status-indicator">
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
                  <div on:click={() => window.open(profileLink(platform))} class="svg-holder">
                    <svg xmlns="http://www.w3.org/2000/svg" height="24" width="24"><path d="M5 21Q4.175 21 3.587 20.413Q3 19.825 3 19V5Q3 4.175 3.587 3.587Q4.175 3 5 3H11Q11.425 3 11.713 3.287Q12 3.575 12 4Q12 4.425 11.713 4.712Q11.425 5 11 5H5Q5 5 5 5Q5 5 5 5V19Q5 19 5 19Q5 19 5 19H19Q19 19 19 19Q19 19 19 19V13Q19 12.575 19.288 12.287Q19.575 12 20 12Q20.425 12 20.712 12.287Q21 12.575 21 13V19Q21 19.825 20.413 20.413Q19.825 21 19 21ZM9 15Q8.725 14.725 8.725 14.3Q8.725 13.875 9 13.6L17.6 5H15Q14.575 5 14.288 4.712Q14 4.425 14 4Q14 3.575 14.288 3.287Q14.575 3 15 3H20Q20.425 3 20.712 3.287Q21 3.575 21 4V9Q21 9.425 20.712 9.712Q20.425 10 20 10Q19.575 10 19.288 9.712Q19 9.425 19 9V6.4L10.375 15.025Q10.1 15.3 9.7 15.3Q9.3 15.3 9 15Z"/></svg>
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
      <div class="setting-label">Caption</div>
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
      <div class="setting-label">Post Delay</div>
      <div class="setting-content post-delay">
        <div class="counter">
          <div on:click={() => changeCounter("PostDelay", "down")} class="minus">
            <svg xmlns="http://www.w3.org/2000/svg" height="24" width="24"><path d="M6 13Q5.575 13 5.287 12.712Q5 12.425 5 12Q5 11.575 5.287 11.287Q5.575 11 6 11H18Q18.425 11 18.712 11.287Q19 11.575 19 12Q19 12.425 18.712 12.712Q18.425 13 18 13Z"/></svg>
          </div>
          <input on:input={e => {filterNonDigits(e, "PostDelay"); $unsavedChanges = true} } value={ instance.PostDelay.num }>
          <div on:click={() => changeCounter("PostDelay", "up") } class="plus">
            <svg xmlns="http://www.w3.org/2000/svg" height="24" width="24"><path d="M12 19Q11.575 19 11.288 18.712Q11 18.425 11 18V13H6Q5.575 13 5.287 12.712Q5 12.425 5 12Q5 11.575 5.287 11.287Q5.575 11 6 11H11V6Q11 5.575 11.288 5.287Q11.575 5 12 5Q12.425 5 12.713 5.287Q13 5.575 13 6V11H18Q18.425 11 18.712 11.287Q19 11.575 19 12Q19 12.425 18.712 12.712Q18.425 13 18 13H13V18Q13 18.425 12.713 18.712Q12.425 19 12 19Z"/></svg>
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
        <div class="setting-label">Queue Delay</div>
        <div class="setting-content queue-delay">
          <div class="counter">
            <div on:click={() => changeCounter("QueueDelay", "down")} class="minus">
              <svg xmlns="http://www.w3.org/2000/svg" height="24" width="24"><path d="M6 13Q5.575 13 5.287 12.712Q5 12.425 5 12Q5 11.575 5.287 11.287Q5.575 11 6 11H18Q18.425 11 18.712 11.287Q19 11.575 19 12Q19 12.425 18.712 12.712Q18.425 13 18 13Z"/></svg>
            </div>
            <input on:input={e => {filterNonDigits(e, "QueueDelay"); $unsavedChanges = true} } style="border-radius: 0px;" value={ instance.QueueDelay.num }>
            <div on:click={() => changeCounter("QueueDelay", "up")} class="plus">
              <svg xmlns="http://www.w3.org/2000/svg" height="24" width="24"><path d="M12 19Q11.575 19 11.288 18.712Q11 18.425 11 18V13H6Q5.575 13 5.287 12.712Q5 12.425 5 12Q5 11.575 5.287 11.287Q5.575 11 6 11H11V6Q11 5.575 11.288 5.287Q11.575 5 12 5Q12.425 5 12.713 5.287Q13 5.575 13 6V11H18Q18.425 11 18.712 11.287Q19 11.575 19 12Q19 12.425 18.712 12.712Q18.425 13 18 13H13V18Q13 18.425 12.713 18.712Q12.425 19 12 19Z"/></svg>
            </div>
          </div>
          <Select on:change={() => $unsavedChanges = true} bind:value={ instance.QueueDelay.unit } width="100px"
            values={{"seconds": "seconds", "minutes": "minutes"}}/>
        </div>
      </div>

      <div class="setting-section">
        <div class="setting-label">Startup Delay</div>
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

</div>
