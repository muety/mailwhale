<script>
  import { onMount } from 'svelte'

  import Layout from '../layouts/Main.svelte'
  import Navigation from '../components/Navigation.svelte'
  import Modal from '../components/Modal.svelte'
  import { getClients, createClient, deleteClient } from '../api/clients'

  const availablePermissions = ['send_mail', 'manage_client', 'manage_template']

  let clients = []

  const emptyClient = {
    id: null,
    description: '',
    api_key: null,
    default_sender: null,
    allowed_senders: null,
    permissions: null,
  }

  let newClientModal
  let newClient
  let newClientPermissions
  let newClientAllowedSenders

  reset()

  async function _createClient() {
    try {
      newClient = await createClient({
        description: newClient.description,
        permissions: Object.entries(newClientPermissions)
          .filter((e) => e[1])
          .map((e) => e[0]),
        default_sender: newClient.default_sender,
        allowed_senders: (newClientAllowedSenders || '').split('\n'),
      })
      clients = [...clients, JSON.parse(JSON.stringify(newClient))]
    } finally {
      newClientModal = false
    }
  }

  async function _deleteClient(id) {
    await deleteClient(id)
    clients = clients.filter((c) => c.id !== id)
    reset()
  }

  function reset() {
    newClient = JSON.parse(JSON.stringify(emptyClient))
    newClientAllowedSenders = null
    newClientPermissions = availablePermissions.reduce(
      (acc, val) => Object.assign(acc, { [val]: false }),
      {}
    )
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
    @apply flex space-x-6 items-center;
  }

  .client-card .info .badges {
    @apply flex space-x-2 mt-1;
  }

  .client-card .info .badges span {
    @apply text-xs bg-primary text-white rounded px-1 font-semibold;
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

      <p class="mb-8">
        <span class="material-icons" style="font-size: inherit;">info</span>
        Clients (aka. API tokens) are used to access MailWhale API from external
        applications, e.g. to send e-mail. Every client is identified by a
        randomly generated ID and a secret, both of which are needed for
        authentication against the API.
      </p>

      {#if newClient.api_key}
        <div
          class="flex flex-col space-y-2 mt-4 mb-12 bg-primary px-4 py-2 w-full rounded text-white text-sm">
          <div>
            <span class="font-semibold">Success!</span>
            <span>A new client was created. Here is your client secret (aka. API
              key). Write it down, as you will not be able to access it later
              on.</span>
          </div>
          <span class="font-mono">{newClient.api_key}</span>
        </div>
      {/if}

      {#if clients.length}
        <div class="flex flex-col space-y-4">
          {#each clients as client, i}
            <div class="client-card">
              <div class="info">
                <span class="text-sm font-semibold">#{i + 1}</span>
                <span
                  class="font-mono text-sm bg-gray-100 p-1 rounded"
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
                  {#if client.default_sender || client.allowed_senders}
                    <div class="flex space-x-2">
                      {#if client.default_sender}
                        <div>
                          <span class="text-xs font-semibold">Default
                            Sender:&nbsp;</span>
                          <span class="text-xs">{client.default_sender}</span>
                        </div>
                      {/if}
                      {#if client.allowed_senders}
                        <div>
                          <span class="text-xs font-semibold">Allowed
                            Senders:&nbsp;</span>
                          <span
                            class="text-xs">{client.allowed_senders.join(', ')}</span>
                        </div>
                      {/if}
                    </div>
                  {/if}
                </div>
              </div>
              <div>
                <a
                  class="text-sm text-primary hover:text-primary-dark underline cursor-pointer"
                  on:click={confirm('Are you sure?') && _deleteClient(client.id)}>Remove</a>
              </div>
            </div>
          {/each}
        </div>
      {:else}
        <div
          class="w-full py-12 text-gray-500 flex justify-center items-center">
          <i>No clients available. Create your first one.</i>
        </div>
      {/if}
    </div>

    {#if newClientModal}
      <Modal on:close={(e) => (newClientModal = false) || reset()}>
        <h1 class="text-2xl font-semibold" slot="header">Add new client</h1>
        <div slot="main" style="min-width: 400px;">
          <form
            class="w-full flex flex-col space-y-4"
            on:submit|preventDefault={_createClient}>
            <div class="flex flex-col w-full space-y-1">
              <label for="desc-input" class="font-semibold">Description</label>
              <input
                type="text"
                class="border-2 border-primary rounded-md p-2"
                name="desc-input"
                placeholder="What will this new client key be used for?"
                required
                bind:value={newClient.description} />
            </div>

            <div>
              <h3 class="font-semibold mb-2 mt-2">Permissions</h3>
              {#each availablePermissions as perm}
                <div class="flex space-x-2 items-center">
                  <input
                    type="checkbox"
                    name="perm-input-{perm}"
                    bind:checked={newClientPermissions[perm]} />
                  <label
                    for="perm-input-{perm}"
                    class="font-mono">{perm}</label>
                </div>
              {/each}
            </div>

            <div>
              <h3 class="font-semibold mb-2 mt-2">E-Mail Settings</h3>
              <div class="flex items-center space-x-2">
                <label
                  for="default-sender-input"
                  class="text-sm font-semibold">Default Sender E-Mail:</label>
                <input
                  type="text"
                  name="default-sender-input"
                  class="border-2 border-primary rounded-md p-2 flex-grow"
                  placeholder="Optional"
                  bind:value={newClient.default_sender} />
              </div>

              <div class="mt-4">
                <label
                  for="allowed-senders-input"
                  class="text-sm font-semibold">Allowed Sender E-Mails:</label>
                <textarea
                  name="allowed-senders-input"
                  class="w-full text-sm border-2 border-primary rounded-md p-2 flex-grow"
                  placeholder="Optional. Line-separated list from allowed e-mail senders (e.g. 'John Doe <john@example.org>')"
                  bind:value={newClientAllowedSenders} />
              </div>
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
