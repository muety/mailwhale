<script>
  import { onMount } from 'svelte'
  import { slide, fly } from 'svelte/transition'

  import Layout from '../layouts/Main.svelte'
  import Navigation from '../components/Navigation.svelte'
  import Modal from '../components/Modal.svelte'
  import { getClients, createClient, deleteClient } from '../api/clients'

  let clients = []

  const emptyClient = { id: null, description: '', apiKey: null }

  let newClientModal
  let newClient

  reset()

  async function _createClient() {
    const responsePayload = await createClient({
      description: newClient.description,
    })
    newClient.id = responsePayload.id
    newClient.description = responsePayload.description
    newClient.apiKey = responsePayload.api_key
    newClientModal = false
	clients = [...clients, JSON.parse(JSON.stringify(newClient))]
  }

  async function _deleteClient(id) {
    await deleteClient(id)
	clients = clients.filter((c) => c.id !== id)
	reset()
  }

  function reset() {
	  newClient = JSON.parse(JSON.stringify(emptyClient))
  }

  onMount(async () => {
    clients = await getClients()
  })
</script>

<style scoped>
  .client-card {
    @apply flex items-center justify-between w-full p-4 border border-gray-300 rounded-md shadow-sm;
  }

  .client-card .info {
    @apply flex space-x-4 items-center;
  }
</style>

<Layout>
  <div slot="content" class="flex">
    <div class="w-1/4">
      <Navigation />
    </div>
    <div class="flex flex-col px-12 w-full w-3/4">
      <div class="flex justify-between mb-8">
        <h1 class="text-2xl font-semibold">Manage API Clients</h1>
        <button
          class="flex items-center px-4 py-2 bg-primary text-white rounded hover:bg-primary-dark"
          on:click|stopPropagation={(e) => (newClientModal = true)}><span
            class="material-icons">add</span>
          Create</button>
      </div>

      <div class="flex items-center mb-8 space-x-2">
        <span class="material-icons">info</span>
        <p>
          Clients (aka. API tokens) are used to access MailWhale API from
          external applications, e.g. to send e-mail. Every client is identified
          by a randomly generated ID and a secret, both of which are needed for
          authentication against the API.
        </p>
      </div>

      {#if newClient.apiKey}
        <div
          class="flex flex-col space-y-2 mt-4 mb-12 bg-primary px-4 py-2 w-full rounded text-white text-sm" transition:fly>
          <div>
            <span class="font-semibold">Success!</span>
            <span>A new client was created. Here is your client secret (aka. API
              key). Write it down, as you will not be able to access it later
              on.</span>
          </div>
          <span class="font-mono">{newClient.apiKey}</span>
        </div>
      {/if}

      <div class="flex flex-col space-y-4">
        {#each clients as client, i}
          <div class="client-card" in:slide>
            <div class="info">
              <span class="text-sm font-semibold">#{i + 1}</span>
              <span
                class="font-mono text-sm bg-gray-100 p-1 rounded"
                title="Client ID">{client.id}</span>
              <span class="text-sm">({client.description})</span>
            </div>
            <div>
              <a
                href="#"
                class="text-sm text-primary hover:text-primary-dark underline"
                on:click={confirm('Are you sure?') && _deleteClient(client.id)}>Remove</a>
            </div>
          </div>
        {/each}
      </div>
    </div>

    {#if newClientModal}
      <Modal on:close={(e) => (newClientModal = false) || reset()}>
        <h1 class="text-2xl font-semibold" slot="header">Add new client</h1>
        <div slot="main" style="min-width: 400px;">
          <form
            class="w-full flex flex-col space-y-4"
            on:submit|preventDefault={_createClient}>
            <div class="flex flex-col w-full space-y-1">
              <label for="desc-input">Description</label>
              <input
                type="desc"
                class="border-2 border-primary rounded-md p-2"
                name="desc-input"
                placeholder="What will this new client key be used for?"
                required
                bind:value={newClient.description} />
            </div>

            <div class="flex justify-between py-2">
              <div />
              <button
                type="submit"
                class="py-2 px-4 text-white bg-primary rounded-md hover:bg-primary-dark">Create</button>
            </div>
          </form>
        </div>
      </Modal>
    {/if}
  </div>
</Layout>
