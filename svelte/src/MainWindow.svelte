<script lang="typescript">
  import Sidebar from './Sidebar.svelte'

  type platform = {
    facebook: string;
    twitter: string;
  }
  type instance = {
    Name: string;
    ImgFolders: string[];
    TimeToQueue: string;
    PostInterval: string;
    PostDelayAtStartup: string;
    Platforms: platform[]; 
  }

  let instances: instance[];

  fetch("./offpost.json")
    .then(response => response.json()).then(data => instances=data);
  console.log("instances: " + instances)
</script>

<style>
  #window-inner {
    border: 5px solid black;
    border-radius: 10px;
    width: 100%;
    height: 700px;
    display: grid;
    grid-template-columns: 100px 1fr;
  }
</style>

<div id="window-inner">
  <Sidebar />
  <div id="mainview">
    <!-- {#each instances as instance}
      <p>{instance.Name}</p>
    {/each} -->
  </div>
</div>
