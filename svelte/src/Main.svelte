<script lang="typescript">
	import MainWindow from './MainWindow.svelte'
  import TwitterConfig from './TwitterConfig.svelte'
	import { activeInstance, ginstances, ginstancesOld, unsavedChanges } from './stores'
  import { fade } from 'svelte/transition'

  // this block is used for frontend testing with 'npm run dev', no backend
  // fetch("test_offpost.json")
  //   .then(response => response.json())
  //   .then(data => $ginstances = toGInst(data) )

  var serversocket = new WebSocket("ws://localhost:14859/config");
  
  // var justOpened: boolean = true
  serversocket.onmessage = async function(e) {
    $ginstances = await toGInst(JSON.parse(e.data))
    $ginstancesOld = JSON.parse(JSON.stringify($ginstances))
  }

  // each message sent over the server socket has an identifier at the beginning
  // of the string in the form of "[idenfifier] "
  // Those identifiers:
  // "s " --> save settings, the message contains the json instance config
  // "d " --> delete instance, message contains the index to delete
  // "twitter " --> configure twitter, message containts index to configure
  // "twitter,r " --> disconnect twitter, message contains index to disconnect
  function sendSocket(msg: string) {
    switch (msg.slice(0, msg.search(" "))) {
      case "twitter,r":
        let i: number = +msg.slice(10)
        delete $ginstances[i].Platforms["twitter"]
        $activeInstance = i-1
        $activeInstance = i
        dialogueActive = false

        $ginstancesOld = JSON.parse(JSON.stringify($ginstances))
        sendSocket("s " + JSON.stringify($ginstances))
        break
      case "twitter":
      case "d":
        serversocket.send(msg)
        break
      case "s":
        let dataBlob: string = msg.slice(2)
        let converted: string = JSON.stringify(fromGInst(JSON.parse(dataBlob)))
        // add the "s, " back so Go can identify the message
        serversocket.send("s " + converted)
        break;
    
      default:
        break;
    }
  }
  
  // returns Go time '10h' as ['10', 'hours']
  function separateTimeUnit(original: string): {num: number, unit: string} {
    let num: string = original.slice(0,-1)
    let unit: string = original.slice(-1)
    switch (unit) {
      case "s":
        unit = "seconds"
        break;
      case "m":
        unit = "minutes"
        break;
      case "h":
        unit = "hours"
        break;
      case "d":
        unit = "days"
        break;
    }
    return {num: +num, unit: unit}
  }

  function dateFromUnixTime(nu: number): string {
    const date = new Date(nu)
    const datearr = date.toString().split(" ")

    const nowDate = new Date()
    const nowDateArr = nowDate.toString().split(" ")

    let fulldate: string = ""
    if (datearr[1] === nowDateArr[1] &&
    +datearr[2]-1 === +nowDateArr[2] &&
    datearr[3] === nowDateArr[3]) {
      fulldate = 'tomorrow '
    } else if (datearr.slice(1,4).every((e, i) => e === nowDateArr.slice(1,4)[i] ) ) {
      fulldate = ""
    } else {
      fulldate = datearr[1] + " " + datearr[2] + ", " + datearr[3] + " "
    }

    let time = datearr[4].split(":")
    let ampm = "am"

    ampm = +time[0] > 11 ? 'pm' : 'am'
    time[0] = +time[0] === 0 ? '12' : time[0]
    time[0] = +time[0] > 12 ? (+time[0] - 12).toString() : time[0]
    time[0] = time[0][0] === '0' ? time[0].slice(-1) : time[0]

    fulldate += "at " + time[0] + ":" + time[1] + ampm

    return fulldate
  }

  async function toGInst(data: any[]): Promise<any[]> {
    for (let i = 0; i < data.length; i++) {
      let e = data[i]
      e.Name = e.Name
      e.ImgFolders = e.ImgFolders
      e.QueueDelay = separateTimeUnit(e.QueueDelay)
      e.PostDelay = separateTimeUnit(e.PostDelay)
      e.StartupPostDelay = e.StartupPostDelay
      // give the UI only the userID of the platform
      e.Platforms = getPlatformUsername(e.Platforms)
      e.Caption = e.Caption
      e.ItemsInQueue = e.ItemsInQueue
      e.NextPostTime = dateFromUnixTime(e.NextPostTime)
      e.Status = e.Status

      let r = await fetch("./userdata/" + e.Name + ".webp")
      if (r.status === 404) {
        e.Image = "./new_instance.svg"
      } else {
        e.Image = "./userdata/" + e.Name + ".webp"
      }
    }
    return data
  }

  function getPlatformUsername(platforms: object): object {
    var resultPlatforms = new Map()
    Object.keys(platforms).forEach(e => {
      if (resultPlatforms.get(e) != "no-config") {
        resultPlatforms.set(e, platforms[e].split("***").slice(-1)[0])
      }
    })
    return Object.fromEntries(resultPlatforms)
  }

  function fromGInst(data: any[]): any[] {
    data.forEach(e => {
      e.Name = e.Name
      e.ImgFolders = e.ImgFolders
      e.QueueDelay = e.QueueDelay["num"] + e.QueueDelay["unit"].slice(0,1)
      e.PostDelay = e.PostDelay["num"] + e.PostDelay["unit"].slice(0,1)
      e.StartupPostDelay = e.StartupPostDelay
      e.Platforms = e.Platforms
      e.Caption = e.Caption
      delete e.ItemsInQueue
      delete e.Image
      delete e.Status
      delete e.NextPostTime
    })
    return data
  }

  function deleteInstance(instance: number) {
    $activeInstance = -1
    $unsavedChanges = false

    dialogueActive = false
    var newGinst = $ginstances.slice(0,instance)
    newGinst.push(...$ginstances.slice(instance+1))

    $ginstances = newGinst

    sendSocket("d " + instance)

  }
  
  var alertText: string
  var alertDisplayed: boolean
  function showAlert(event) {
    alertText = event.detail.text

    if (!alertDisplayed) {
      alertDisplayed = true
      setTimeout(() => {
        alertDisplayed = false
      }, 3000);
    }
  }  

  var dialogueActive: boolean = false
  var dialogueType: string
  var instIndex: number = 0
  function showDialogue(event) {
    dialogueType = event.detail.type
    dialogueActive = true
    if (instIndex > -1) {
      instIndex = event.detail.instance
    }
  }
  
</script>

<style>
	#the-container {
    height: 100vh;
		margin: auto;
		display: grid;
		justify-items: center;
    align-items: start;
    overflow: hidden;
    position: relative;
	}

  #mainwindow-wrapper {
    width: 800px;
		display: grid;
		place-items: center;
  }

  #brand-header {
    display: grid;
    align-items: end;
    grid-template-columns: auto auto;
    width: 100%;
    gap: 10px;
    margin-bottom: 9px;
    height: 100px;
  }

  #logo-text {
    display: flex;
    align-items: center;
    gap: 3px;
  }

  #logo {
    width: 40px;
  }

  #alert {
    white-space: nowrap;
    overflow: hidden;
    border-radius: 10px;
    margin: 30px 0px 30px 0px;
    background: #c7c7c7;
    text-align: center;
    opacity: 0;
    z-index: 2;
    position: absolute;
    top: 29px;

    transition: opacity 0.2s;
  }

  #alert.alertDisplayed {
    opacity: 1;
    transition: opacity 0.2s;
  }

  #alert span {
    display: block;
    margin: 10px 10px 10px 10px;
  }
  
  #dialogue {
    position: absolute;
    width: 100%;
    height: 100%;
    background: rgba(0,0,0,0.2);
    display: grid;
    align-items: start;
    justify-items: center;
    padding-top: 300px;
    z-index: 10;
  }

  @media(max-height: 900px) {
    #dialogue {
    align-items: center;
    padding-top: 0px;
    }
  }

  #dialogue-inner {
    background: white;
    border-radius: 10px;
    /* box-shadow: 0px 0px 20px -10px black; */
    box-sizing: border-box;
  }
  #delete-content {
    display: grid;
    grid-template-rows: auto 30px;
    padding: 20px;
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

<div id="the-container">
  <div class:alertDisplayed id="alert">
     <span> {alertText} </span>
  </div>

  {#if dialogueActive}
    <div transition:fade="{{ duration: 70 }}" id="dialogue" on:click={() => {dialogueActive = false}}>
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
        {:else if dialogueType === "twitter"}
          <TwitterConfig instance={instIndex} 
            on:socketMessage={e => sendSocket(e.detail.text)}
          />
        {:else}
          <div>
            we dont have a dialogue setup for this yet.
          </div>
        {/if}
      </div>
    </div>
  {/if}

  <div id="mainwindow-wrapper">
    <div id="brand-header">
      <div id="logo-text">
        <img id="logo" src="./logo.svg" alt="Offpost logo">
        <span>Offpost</span>
      </div>
    </div>
    <MainWindow
      on:dialogue={showDialogue}
      on:alert={showAlert}
      on:socketMessage={m => sendSocket(m.detail.text)}
    />
  </div>
</div>
