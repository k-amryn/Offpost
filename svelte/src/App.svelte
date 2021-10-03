<script lang="typescript">
	import MainWindow from './MainWindow.svelte';
	import { ginstances } from './stores'

  var serversocket = new WebSocket("ws://localhost:8081/config");

  serversocket.onmessage = function(e) {
      $ginstances = toGInst(JSON.parse(e.data))
  };
  
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


  function toGInst(data: any[]): any[] {
    data.forEach(e => {
      e.Name = e.Name
      e.ImgFolders = e.ImgFolders
      e.TimeToQueue = separateTimeUnit(e.TimeToQueue)
      e.PostInterval = separateTimeUnit(e.PostInterval)
      e.PostDelayAtStartup = e.PostDelayAtStartup
      e.Platforms = e.Platforms
      e.ItemsInQueue = e.ItemsInQueue
      e.NextPostTime = dateFromUnixTime(e.NextPostTime)
      e.Status = e.Status
    })
    return data
  }

</script>

<style>
	#the-container {
		width: 800px;
		margin: auto;
		display: grid;
		place-items: center;
	}
</style>

<div id="the-container">
	<div id="brand-heading">
		<img src="./logo.svg" alt="Offpost logo" width="200px">
	</div>
	<MainWindow />
</div>
