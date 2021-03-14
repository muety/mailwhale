<script>
  import { slide } from 'svelte/transition'
  import AccountIndicator from '../components/AccountIndicator.svelte'

  import { user } from '../stores/auth'
  import { errors, infos, successes } from '../stores/alerts'

  function logout() {
    user.logout()
    window.location.replace('/')
  }
</script>

<div id="app-container" class="container mx-auto my-8 flex flex-col flex-grow">
  <div
    id="alert-container"
    class="w-full absolute inset-x-0 top-0 py-8 flex justify-center space-y-2 flex-col items-center">
    {#each $errors as m}
      <div
        class="flex space-x-2 mt-4 bg-red-500 px-4 py-2 rounded text-white text-sm"
        transition:slide>
        <span class="material-icons">warning</span>
        <span>{m}</span>
      </div>
    {/each}
    {#each $infos as m}
      <div
        class="flex space-x-2 mt-4 bg-primary px-4 py-2 rounded text-white text-sm"
        transition:slide>
        <span class="material-icons">info</span>
        <span>{m}</span>
      </div>
    {/each}
    {#each $successes as e}
      <div
        class="flex space-x-2 mt-4 bg-green-500 px-4 py-2 rounded text-white text-sm"
        transition:slide>
        <span class="material-icons">check_circle</span>
        <span>{e}</span>
      </div>
    {/each}
  </div>

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
      <a href="imprint" class="text-primary hover:text-primary-dark">Imprint &
        Data Privacy</a>
    </div>
  </footer>
</div>
