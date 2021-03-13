<script>
  import router from 'page'
  import Home from './views/Home.svelte'
  import Clients from './views/Clients.svelte'

  import { user } from './stores/auth'

  let page

  user.load()
  user.subscribe((user) => {
    if (!user) {
      page = Home
    }
  })

  router('/', () => (page = Home))
  router('/login', () => (page = Home))
  router('/clients', () => (page = Clients))

  router.start()
</script>

<div id="app-container" class="container mx-auto flex flex-col flex-grow">
  <svelte:component this={page} />
</div>
