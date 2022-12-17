<script>
  import { createEventDispatcher } from 'svelte'

  const dispatch = createEventDispatcher()

  function deleteClient(id) {
    dispatch('delete', { id })
  }

  export let client
</script>

<style scoped>
  .client-card {
    @apply flex items-center justify-between w-full p-4 border border-gray-300 rounded-md shadow-sm;
  }

  .client-card .info {
    @apply flex space-x-6 items-center;
  }

  .client-card .info .badges {
    @apply flex space-x-2 mt-1;
  }

  .client-card .info .badges span {
    @apply text-xs bg-primary text-white rounded px-1 font-semibold;
  }

  .client-card .actions .badges span {
    @apply text-sm bg-primary text-white rounded p-2 font-semibold;
  }
</style>

<div class="client-card">
  <div class="info">
    <span class="text-sm font-semibold">
      #<slot name="index"></slot>
    </span>
    <span
      class="p-1 font-mono text-sm bg-gray-100 rounded"
      title="Client ID">{client.id}</span>
    <div class="flex flex-col">
      <span class="text-sm">{client.description}</span>
      {#if client.permissions}
        <div class="badges">
          {#each client.permissions as perm}
            <span>{perm}</span>
          {/each}
        </div>
      {/if}
      {#if client.sender}
        <div class="flex space-x-2">
          <div>
            <span class="text-xs font-semibold">Sender E-Mail:&nbsp;</span>
            <span class="text-xs">{client.sender}</span>
          </div>
        </div>
      {/if}
    </div>
  </div>
  <div class="flex space-x-4 actions">
    <div class="badges"><span title="{client.count_mails} mails sent">{client.count_mails}</span></div>
    <span
      class="text-sm underline cursor-pointer text-primary hover:text-primary-dark"
      on:click={confirm('Are you sure?') && deleteClient(client.id)}>Remove</span>
  </div>
</div>
