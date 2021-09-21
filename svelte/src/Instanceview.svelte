<script lang="typescript">
  import { ginstances, activeInstance } from './stores'

  let advanced: boolean = false
  $: instance = $ginstances[$activeInstance]

  function addFolder() {
    instance.ImgFolders = [...instance.ImgFolders, ""]
  }

  function removeFolder(i: number) {
    instance.ImgFolders.splice(i, 1)
    instance.ImgFolders = instance.ImgFolders
  }

  function addVariable(variable: string) {
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

  // used to ensure number inputs are only numbers
  function filterNonDigits(e: any) {
    e.target.value = e.target.value.replace(/\D/, '')
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
    margin-top: 0px;
    cursor: pointer;
    display: inline-block;
    transform: translateY(-7px)
  }

  .setting-section {
    margin-top: 20px;
    display: grid;
    grid-template-columns: 115px 1fr;
  }

  .setting-label {
    margin-top: 7px;
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
    padding-top: 0px;
  }

  .svg-holder:not(.status-indicator) {
    cursor: pointer;
  }

  .folder-input {
    width: 100%;
  }

  select, button, input {
    background: transparent;
    border: 2px solid black;
    border-radius: 5px;
    height: 36px;
    padding-top: 0px;
    padding-bottom: 0px;
  }

  select, button {
    cursor: pointer;
  }

  .platform-row {
    display: grid;
    grid-template-columns: min-content auto auto auto;
    width: min-content;
    gap: 10px;
  }

  .platform-row .svg-holder {
    margin-left: 10px;
  }

  .caption-input {
    width: 100%;
    border: 2px solid black;
    border-radius: 5px;
    resize: vertical;
  }

  .caption-subtext {
    font-size: var(--small-font);
    display: inline-block;
    transform: translateY(-9px)
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
    gap: 10px;
  }

  .minus, .plus {
    height: 32px;
    margin-bottom: 7px;
    display: grid;
    place-items: center;
    cursor: pointer;
  }

  .minus {
    border-top: 2px solid black;
    border-bottom: 2px solid black;
    border-left: 2px solid black;
    border-top-left-radius: 5px;
    border-bottom-left-radius: 5px;
  }

  .plus {
    border-top: 2px solid black;
    border-bottom: 2px solid black;
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
    gap: 6px;
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
    /* background: red; */
  }

  #status-bar span {
    font-size: var(--small-font)
  }

  #status-bar button {
    margin-top: 6px;
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
</style>

<div id="container">
  <div id="settings">
    <div id="heading-bar">
      <p id="heading-text">{ instance.Name }</p>
    </div>
    <div class="setting-section">
      <div class="setting-label">Folders:</div>
      <div class="setting-content">
        {#each instance.ImgFolders as folder, i}
        <div class="folder-row">
            <input class="folder-input" bind:value={folder}>
            <div class="svg-holder">
              <svg width="18px" style="margin-bottom: 12px;" viewBox="0 0 25 21" xmlns="http://www.w3.org/2000/svg" fill-rule="evenodd" clip-rule="evenodd" stroke-linecap="round" stroke-linejoin="round" stroke-miterlimit="1.5">
                <path d="M2412 271v31h47v-26h-25l-4-5h-18z" fill="none" stroke="#000" stroke-width="5.8" transform="matrix(.45762 0 0 .56483 -1102 -151)"/>
              </svg>
            </div>
            <div class="svg-holder" on:click={() => removeFolder(i)}>
              <svg width="12px" style="margin-bottom: 10px;" version="1.1" viewBox="0 0 5.4152 5.4152" xmlns="http://www.w3.org/2000/svg">
                <g transform="translate(-322.81 -103.89)" fill="none" stroke="#de3e39" stroke-linecap="round">
                <path d="m323.31 104.39 4.4152 4.4152"/>
                <path d="m327.73 104.39-4.4152 4.4152"/>
                </g>
              </svg>
            </div>
          </div>
        {/each}
        <span on:click={() => addFolder()} class="subtext">Add new folder</span>

      </div>
    </div>

    <div class="setting-section">
      <div class="setting-label">Platforms:</div>
      <div class="setting-content">
        {#each Object.keys(instance.Platforms) as platform}
          <div class="platform-row">
            <select>
              <option>Twitter</option>
              <option>Facebook</option>
              <option>Tumblr</option>
            </select>
            <button>Configure</button>
            <div class="svg-holder status-indicator">
              <svg width="12px" style="margin-bottom: 8px" viewBox="0 0 15 15" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" xml:space="preserve">
                <g transform="matrix(1,0,0,1,-1223.59,-1560.19)">
                <g transform="matrix(1,0,0,1,0,1101.51)">
                <g transform="matrix(0.521739,0,0,0.521739,591.235,-354.719)">
                <circle cx="1225.71" cy="1572.7" r="13.69" style="fill:rgb(65,147,62);"/>
                </g>
                </g>
                </g>
              </svg>
            </div>
            <div class="svg-holder">
              <svg width="14px" style="margin-bottom: 10px;fill-rule:evenodd;clip-rule:evenodd;stroke-linecap:round;stroke-linejoin:round;stroke-miterlimit:1.5;" viewBox="0 0 15 15" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" xml:space="preserve">
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
          </div>
        {/each}
        <span class="subtext">Add new platform</span>
      </div>
    </div>

    <div class="setting-section">
      <div class="setting-label">Caption:</div>
      <div class="setting-content">
        <textarea class="caption-input" rows="2"></textarea>
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
          <div on:click={() => {if (instance.PostInterval.num > 0) instance.PostInterval.num-=1} } class="minus">
            <svg width="15px" version="1.1" viewBox="0 0 11.863 11.863" xmlns="http://www.w3.org/2000/svg">
              <g transform="translate(-173.99 -176.63)" fill="none" stroke="#000" stroke-linecap="round" stroke-width="1.5218px">
                <path d="m185.09 182.56h-10.341"/>
              </g>
            </svg>
          </div>
          <input on:input={e => filterNonDigits(e) } bind:value={ instance.PostInterval.num }>
          <div on:click={() => instance.PostInterval.num+=1 } class="plus">
            <svg width="15px" version="1.1" viewBox="0 0 11.863 11.863" xmlns="http://www.w3.org/2000/svg">
              <g transform="translate(-173.99 -176.63)" fill="none" stroke="#000" stroke-linecap="round" stroke-width="1.5218px">
                <path d="m179.92 177.39v10.341"/>
                <path d="m185.09 182.56h-10.341"/>
              </g>
            </svg>
          </div>
        </div>
        <select bind:value={ instance.PostInterval.unit }>
          <option>minutes</option>
          <option>hours</option>
          <option>days</option>
        </select>
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
            <div on:click={() => {if (instance.TimeToQueue.num > 0) instance.TimeToQueue.num-=1} } class="minus">
              <svg width="15px" version="1.1" viewBox="0 0 11.863 11.863" xmlns="http://www.w3.org/2000/svg">
                <g transform="translate(-173.99 -176.63)" fill="none" stroke="#000" stroke-linecap="round" stroke-width="1.5218px">
                  <path d="m185.09 182.56h-10.341"/>
                </g>
              </svg>
            </div>
            <input on:input={e => filterNonDigits(e) } style="border-radius: 0px;" bind:value={instance.TimeToQueue.num}>
            <div on:click={() => instance.TimeToQueue.num+=1 } class="plus">
              <svg width="15px" version="1.1" viewBox="0 0 11.863 11.863" xmlns="http://www.w3.org/2000/svg">
                <g transform="translate(-173.99 -176.63)" fill="none" stroke="#000" stroke-linecap="round" stroke-width="1.5218px">
                  <path d="m179.92 177.39v10.341"/>
                  <path d="m185.09 182.56h-10.341"/>
                </g>
              </svg>
            </div>
          </div>
          <select>
            <option>seconds</option>
            <option>minutes</option>
          </select>
        </div>
      </div>

      <div class="setting-section">
        <div class="setting-label">Startup Delay:</div>
        <div class="setting-content">
          <select bind:value={ instance.PostDelayAtStartup }>
            <!-- On startup, attempt posting at the NextPostTime. If the NextPostTime has passed, 
            then use this option. -->
            <option value="random">Random</option>
            <option value="full">Full</option>
            <option value="none">None</option>
          </select>
        </div>
      </div>

      <div class="delete-instance">
        <span>Delete instance</span>
      </div>
    {/if}
  </div>
  <div id="status-bar">
    {#if instance.Status === "needs-configuring"}
      <span class="needs-configuring">Needs configuring</span>
    {:else if instance.ItemsInQueue > 0}
      <span class="queued">{instance.ItemsInQueue} items in queue<br>
        Next post at {instance.NextPostTime}<br></span>
    {:else if instance.ItemsInQueue === 0}
      <span class="waiting">Waiting for new image</span>
    {/if}
    <div id="status-buttons">
      <button>Pause</button>
    </div>
  </div>

</div>
