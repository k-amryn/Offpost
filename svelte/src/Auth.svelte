<script lang="typescript">
  var serversocket = new WebSocket("ws://localhost:14859/authdata");
  
  var platform: string
  var instName: string
  var connected: boolean
  serversocket.onopen = () => {
    serversocket.send(window.location.search)
  }
    
  serversocket.onmessage =  e => {
    let info: string = e.data
    let i = info.search(" ")
    platform = info.slice(0, i)
    instName = info.slice(i+1)
     
    if (window.location.search.slice(0,6) === "?error") {
      connected = false
    } else {
      connected = true
    }
    serversocket.close()
  }
</script>

<style>
  #wrap {
    width: 100vw;
    height: 100vh;
    display: grid;
    place-items: center;
    justify-items: center;
  }
  #info {
    display: flex;
    gap: 10px;
    height: 200px;
    flex-direction: column;
  }

  #logo {
    height: 40px;
    gap: 3px;
    display: flex;
    align-items: center;
  }

  #logo img {
    width: 40px;
  }

  #alert {
    width: 400px;
    height: 110px;
    overflow: hidden;
    border-radius: 10px;
    background: white;
    padding: 15px;
  }
</style>

<div id="wrap">
  <div id="info">
    <div id="logo">
      <img src="./logo.svg" alt="Offpost logo">
      Offpost
    </div>
    <div id="alert">
      {#if connected}
      <span>{instName}: Connected to {platform}.</span>
      {:else}
      <span>{instName}: {platform} connection cancelled.</span>
      {/if}
    </div>
  </div>
</div>