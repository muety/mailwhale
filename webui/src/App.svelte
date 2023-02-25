<script>
  import router from 'page'
  import Home from './views/Home.svelte'
  import Login from './views/Login.svelte'
  import Signup from './views/Signup.svelte'
  import Clients from './views/Clients.svelte'
  import Mails from './views/Mails.svelte';
  import Imprint from './views/Imprint.svelte';
  import Addresses from './views/Addresses.svelte';
  import Templates from './views/Templates.svelte';
  import { sanitize } from './utils/url'
  import { config } from './stores/config'

  import { user } from './stores/auth'

  let page

  user.load()

  const basePath = config.get().basePath

  router(`/${basePath}`, () => {
    if (!!user.getToken()) {
      router.redirect(sanitize(`/${basePath}/clients`))
    }
    page = Home
  })
  router(sanitize(`/${basePath}/login`), () => (page = Login))
  router(sanitize(`/${basePath}/signup`), () => (page = Signup))
  router(sanitize(`/${basePath}/clients`), () => (page = Clients))
  router(sanitize(`/${basePath}/mails`), () => (page = Mails))
  router(sanitize(`/${basePath}/imprint`), () => (page = Imprint))
  router(sanitize(`/${basePath}/addresses`), () => (page = Addresses))
  router(sanitize(`/${basePath}/templates`), () => (page = Templates))

  // don't forget to update stores/config/reservedRoutes when adding a new route here

  router.start()
</script>

<style global lang="postcss">
  :root {
    --color-primary: #159ce4;
    --color-primary-dark: #138dce;
    --color-primary-light: #7dc0e4;
    --color-text: #4b5563;
  }

  @tailwind base;
  @tailwind components;
  @tailwind utilities;
</style>

<div id="app-container" class="container flex flex-col flex-grow mx-auto">
  <svelte:component this={page}/>
</div>
