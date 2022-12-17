<script>
  import { createEventDispatcher } from 'svelte'

  const dispatch = createEventDispatcher()

  function logout() {
    dispatch('logout')
  }

  export let currentUser
</script>

<div class="flex items-center h-full space-x-2">
  {#if currentUser}
    <span class="flex items-center text-sm"><span class="mr-1 material-icons">person</span> {currentUser}</span>
    <span>|</span>
    <a
      class="text-sm cursor-pointer text-primary hover:text-primary-dark"
      on:click={logout}>Logout</a>
  {:else}
    <a
      href="/signup"
      class="font-semibold text-primary hover:text-primary-dark">Sign Up</a>
    <span>&nbsp;|&nbsp;</span>
    <a
      href="/login"
      class="font-semibold text-primary hover:text-primary-dark">Login</a>
  {/if}
</div>
