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

  function toGInst(data: any[]): any[] {
    data.forEach(e => {
      e.Name = e.Name
      e.ImgFolders = e.ImgFolders
      e.TimeToQueue = separateTimeUnit(e.TimeToQueue)
      e.PostInterval = separateTimeUnit(e.PostInterval)
      e.PostDelayAtStartup = e.PostDelayAtStartup
      e.Platforms = e.Platforms
      e.ItemsInQueue = e.ItemsInQueue
      e.NextPostTime = e.NextPostTime
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
