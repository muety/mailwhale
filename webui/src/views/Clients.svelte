<script>
  import {onMount} from 'svelte'

  import Layout from '../layouts/Main.svelte'
  import Navigation from '../components/Navigation.svelte'
  import Modal from '../components/Modal.svelte'
  import {createClient, deleteClient, getClients} from '../api/clients'
  import {getMe} from '../api/users'
  import ClientCard from '../components/ClientCard.svelte'

  const availablePermissions = [
    'send_mail',
    'manage_client',
    'manage_user',
    'manage_template'
  ]

  let me
  let clients = []

  const emptyClient = {
    id: null,
    description: '',
    api_key: null,
    sender: null,
    permissions: null,
    count_mails: 0
  }

  let newClientModal
  let newClient
  let newClientPermissions

  reset()

  async function _createClient() {
    try {
      newClient = await createClient({
        description: newClient.description,
        permissions: Object.entries(newClientPermissions)
          .filter((e) => e[1])
          .map((e) => e[0]),
        sender: newClient.sender
      })
      clients = [...clients, JSON.parse(JSON.stringify(newClient))]
    } finally {
      newClientModal = false
    }
  }

  async function _deleteClient({ detail }) {
    await deleteClient(detail.id)
    clients = clients.filter((c) => c.id !== detail.id)
    reset()
  }

  function reset() {
    newClient = JSON.parse(JSON.stringify(emptyClient))
    newClientPermissions = availablePermissions.reduce(
      (acc, val) => Object.assign(acc, {[val]: false}),
      {}
    )
  }

  onMount(async () => {
    me = await getMe()
    clients = await getClients()
  })
</script>

<Layout>
  <div slot="content" class="flex">
    <div class="w-1/4">
      <Navigation/>
    </div>
    <div class="flex flex-col w-3/4 w-full px-12">
      <div class="flex justify-between mb-8">
        <h1 class="text-2xl font-semibold">Manage API Clients</h1>
        <button
          class="flex items-center px-4 py-2 text-white rounded bg-primary hover:bg-primary-dark"
          on:click|stopPropagation={(e) => (newClientModal = true) && reset()}><span
          class="material-icons">add</span>
          Create
        </button>
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
          class="flex flex-col w-full px-4 py-2 mt-4 mb-12 space-y-2 text-sm text-white rounded bg-primary">
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
            <ClientCard client="{client}" on:delete={_deleteClient}>
              <span slot="index">{i+1}</span>
            </ClientCard>
          {/each}
        </div>
      {:else}
        <div
          class="flex items-center justify-center w-full py-12 text-gray-500">
          <i>No clients available. Create your first one.</i>
        </div>
      {/if}
    </div>

    {#if newClientModal}
      <Modal on:close={(e) => (newClientModal = false) || reset()}>
        <h1 class="text-2xl font-semibold" slot="header">Add new client</h1>
        <div slot="main" style="min-width: 400px;">
          <form
            class="flex flex-col w-full space-y-4"
            on:submit|preventDefault={_createClient}>
            <div class="flex flex-col w-full space-y-1">
              <label for="desc-input" class="font-semibold">Description</label>
              <input
                type="text"
                class="p-2 border-2 rounded-md border-primary"
                id="desc-input"
                name="desc-input"
                placeholder="What will this new client key be used for?"
                required
                bind:value={newClient.description}/>
            </div>

            <div>
              <h3 class="mt-2 mb-2 font-semibold">Permissions</h3>
              {#each availablePermissions as perm}
                <div class="flex items-center space-x-2">
                  <input
                    type="checkbox"
                    id="perm-input-{perm}"
                    name="perm-input-{perm}"
                    bind:checked={newClientPermissions[perm]}/>
                  <label
                    for="perm-input-{perm}"
                    class="font-mono">{perm}</label>
                </div>
              {/each}
            </div>

            <div>
              <h3 class="mt-2 mb-2 font-semibold">E-Mail Settings</h3>
              <div class="max-w-screen-md mb-4">
                <div
                  class="px-4 py-2 mt-4 text-sm text-white rounded bg-primary">
                  <span class="font-semibold">Please Note:</span>
                  <span>You can set an optional sender address for this client
                    (e.g.
                    <strong><i>My App &lt;noreply@example.org&gt;</i></strong>),
                    that will be used in the mail's
                    <i>"From"</i>
                    header. However, you need to make sure that SPF and DMARC
                    records are properly set for your domain. You need to
                    authorize MailWhale's servers to send mail on your behalf.
                    If left blank, a default sender address like
                    <strong><i>vldsbgfr+user@mailwhale.dev</i></strong></span>
                  will be used.
                </div>
              </div>
              <div class="flex items-center space-x-2">
                <label
                  for="default-sender-input"
                  class="text-sm font-semibold">Sender E-Mail:</label>
                <select
                  class="border-2 border-primary rounded-md p-2 flex-grow cursor-pointer"
                  id="default-sender-input"
                  bind:value={newClient.sender}>
                  <option selected value>Default</option>
                  {#each me.senders as sender}
                    <option value={sender.mail}>
                      {sender.mail}
                      ({sender.verified ? 'verified' : 'not verified'})
                    </option>
                  {/each}
                </select>
              </div>
            </div>

            <div class="flex justify-between py-2">
              <div/>
              <button
                type="submit"
                class="px-4 py-2 text-white rounded-md bg-primary hover:bg-primary-dark">Create
              </button>
            </div>
          </form>
        </div>
      </Modal>
    {/if}
  </div>
</Layout>
