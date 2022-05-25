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

  import { user } from './stores/auth'

  let page

  user.load()

  router('/', () => {
    if (!!user.getToken()) {
      router.redirect('/clients')
    }
    page = Home
  })
  router('/login', () => (page = Login))
  router('/signup', () => (page = Signup))
  router('/clients', () => (page = Clients))
  router('/mails', () => (page = Mails))
  router('/imprint', () => (page = Imprint))
  router('/addresses', () => (page = Addresses))
  router('/templates', () => (page = Templates))

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
  <svelte:component this={page} />
</div>
