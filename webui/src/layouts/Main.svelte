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

<div class="container flex flex-col flex-grow mx-auto my-8">
  <div
    id="alert-container"
    class="absolute inset-x-0 top-0 flex flex-col items-center justify-center w-full py-8 space-y-2">
    {#each $errors as m}
      <div
        class="flex px-4 py-2 mt-4 space-x-2 text-sm text-white bg-red-500 rounded"
        transition:slide>
        <span class="material-icons">warning</span>
        <span>{m}</span>
      </div>
    {/each}
    {#each $infos as m}
      <div
        class="flex px-4 py-2 mt-4 space-x-2 text-sm text-white rounded bg-primary"
        transition:slide>
        <span class="material-icons">info</span>
        <span>{m}</span>
      </div>
    {/each}
    {#each $successes as e}
      <div
        class="flex px-4 py-2 mt-4 space-x-2 text-sm text-white bg-green-500 rounded"
        transition:slide>
        <span class="material-icons">check_circle</span>
        <span>{e}</span>
      </div>
    {/each}
  </div>

  <header class="flex justify-between w-full">
    <a id="logo-container" href="/" class="flex items-center space-x-4">
      <img src="images/logo.svg" alt="Logo" style="max-height: 60px;"/>
      <div class="flex">
        <span class="text-xl font-semibold text-primary">Mail</span>
        <span class="text-xl font-semibold">Whale</span>
      </div>
    </a>
    <div>
      <AccountIndicator currentUser={$user} on:logout={logout}/>
    </div>
  </header>

  <main class="flex-grow mt-24">
    <slot name="content"/>
  </main>

  <footer class="flex justify-between mt-8 text-sm">
    <div class="flex space-x-4">
      <a
        href="https://github.com/muety/mailwhale"
        target="_blank"
        class="text-primary hover:text-primary-dark">GitHub</a>
    </div>
    <div class="flex space-x-4">
      <a href="imprint" class="text-primary hover:text-primary-dark">Imprint &
        Data Privacy</a>
    </div>
  </footer>
</div>
