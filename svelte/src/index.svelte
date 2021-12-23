<script lang="typescript">
  import { Router, Route } from 'svelte-routing'
  import Main from './Main.svelte'
  import Auth from './Auth.svelte';

  export let url = ''
</script>

<Router {url}>
  <Route path="/" component={Main} />
  <Route path="/auth" component={Auth} />
</Router>