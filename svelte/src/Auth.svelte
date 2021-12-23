<script lang="typescript">
  var serversocket = new WebSocket("ws://localhost:14859/authdata");
  
  var platform: string
  var instName: string
  serversocket.onopen = () => {
    serversocket.send(window.location.search)
  }
    
  serversocket.onmessage =  e => {
    let info: string = e.data
    let i = info.search(" ")
    platform = info.slice(0, i)
    instName = info.slice(i+1)
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
    display: grid;
    height: 200px;
    grid-template-columns: auto auto;
    place-content: center;
  }

  #logo {
    width: 200px;
  }

  #alert {
    width: 400px;
    height: 110px;
    white-space: nowrap;
    overflow: hidden;
    border-radius: 10px;
    border: var(--main-border-size);
    margin: 30px 0px 30px 0px;
    padding: 10px;
  }
</style>

<div id="wrap">
  <div id="info">
    <img id="logo" src="./logo.svg" alt="Offpost logo">
    <div id="alert">
      <span>{instName}: Connected to {platform}.</span>
    </div>
  </div>
</div>