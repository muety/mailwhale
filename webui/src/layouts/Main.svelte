<script>
  import AccountIndicator from '../components/AccountIndicator.svelte'

  import { user } from '../stores/auth'

  function logout() {
    user.logout()
    window.location.replace('/')
  }
</script>

<div id="app-container" class="container mx-auto my-8 flex flex-col flex-grow">
  <header class="flex w-full justify-between">
    <div id="logo-container" class="flex space-x-4 items-center">
      <img src="images/logo.svg" alt="Logo" style="max-height: 60px;" />
      <div class="flex">
        <span class="text-primary text-xl font-semibold">Mail</span>
        <span class="text-xl font-semibold">Whale</span>
      </div>
    </div>
    <div>
      <AccountIndicator currentUser={$user} on:logout={logout} />
    </div>
  </header>

  <main class="mt-24 flex-grow">
    <slot name="content" />
  </main>

  <footer class="flex justify-between text-sm">
    <div class="flex space-x-4">
      <a
        href="https://github.com/muety/mailwhale"
        class="text-primary hover:text-primary-dark">GitHub</a>
    </div>
    <div class="flex space-x-4">
      <a href="#" class="text-primary hover:text-primary-dark">Imprint & Data
        Privacy</a>
    </div>
  </footer>
</div>
