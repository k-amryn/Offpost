<script lang="typescript">
	import MainWindow from './MainWindow.svelte'
	import { activeInstance, ginstances, ginstancesOld } from './stores'

  // this block is used for frontend testing with 'npm run dev', no backend
  // fetch("test_offpost.json")
  //   .then(response => response.json())
  //   .then(data => {
  //     $ginstances = toGInst(data)
  //   })

  var serversocket = new WebSocket("ws://localhost:14859/config");
  
  // var justOpened: boolean = true
  serversocket.onmessage = async function(e) {
    $ginstances = await toGInst(JSON.parse(e.data))
    $ginstancesOld = JSON.parse(JSON.stringify($ginstances))
  }

  let dialogueActive: boolean = false

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
      e.Platforms = getPlatformUserID(e.Platforms)
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

  function getPlatformUserID(platforms: object): object {
    var resultPlatforms = new Map()
    Object.keys(platforms).forEach(e => {
      resultPlatforms.set(e, platforms[e].split("***").slice(-1)[0])
      
      if (resultPlatforms.get(e) != "no" && resultPlatforms.get(e) != "no-config") {
        switch (e) {
          case "facebook":
            break
          case "twitter":
            resultPlatforms.set(e, "https://twitter.com/i/user/" + resultPlatforms.get(e))
            break

        }
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
</script>

<style>
	#the-container {
		width: 800px;
		margin: auto;
		display: grid;
		place-items: center;
	}
  #brand-heading {
    display: grid;
    width: 90%;
    height: 200px;
    grid-template-columns: auto auto;
    place-content: center;
    align-content: center;
  }

  #logo {
    width: 200px;
  }

  #alert {
    width: 0px;
    height: 130px;
    white-space: nowrap;
    overflow: hidden;
    border-radius: 10px;
    border: var(--main-border-size);
    margin: 30px 0px 30px 0px;
    opacity: 0;

    transition: width 0.3s, opacity 0.3s;
  }

  #alert.alertDisplayed {
    width: 400px;
    opacity: 1;
    transition: width 0.3s, opacity 0.3s;
  }

  #alert span {
    display: block;
    margin: 10px 10px 10px 10px;
  }
</style>

<div id="the-container">
	<div id="brand-heading">
		<img id="logo" src="./logo.svg" alt="Offpost logo">
    <div class:alertDisplayed id="alert">
       <span> {alertText} </span>
    </div>
	</div>
	<MainWindow
    bind:dialogueActive={dialogueActive}
    on:alert={showAlert}
    on:socketMessage={m => sendSocket(m.detail.text)}
  />
</div>
