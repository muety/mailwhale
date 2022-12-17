<script>
  import { onMount } from 'svelte'

  import Layout from '../layouts/Main.svelte'
  import Navigation from '../components/Navigation.svelte'
  import Modal from '../components/Modal.svelte'
  import { getMe, updateMe } from '../api/users'
  import { successes } from '../stores/alerts'

  let me
  let newAddress = ''
  let newAddressModal = false

  async function _deleteAddress(sender) {
    me = await updateMe({
      senders: me.senders.filter((s) => s.mail !== sender.mail).map(s => s.mail),
    })
    successes.spawn('Sender addresses removed')
  }

  async function _addAddress() {
    try {
      me = await updateMe({
        senders: [...me.senders.map((s) => s.mail), newAddress],
      })
      successes.spawn(
        'Sender addresses added. Please check your inbox for verification.'
      )
    } finally {
      newAddressModal = false
      reset()
    }
  }

  function reset() {
    newAddress = ''
  }

  onMount(async () => {
    me = await getMe()
  })
</script>

<style scoped>
  .address-card {
    @apply flex items-center justify-between w-full p-4 border border-gray-300 rounded-md shadow-sm;
  }

  .address-card .info {
    @apply flex space-x-6 items-center;
  }

  .address-card .info .badges {
    @apply flex space-x-2 mt-1;
  }

  .address-card .info .badges span {
    @apply text-xs text-white rounded px-1 font-semibold;
  }
</style>

<Layout>
  <div slot="content" class="flex">
    <div class="w-1/4">
      <Navigation/>
    </div>
    <div class="max-w-screen-lg">
      <div class="flex justify-between mb-8">
        <h1 class="text-2xl font-semibold">E-Mail Addresses</h1>
        <button
          class="flex items-center px-4 py-2 text-white rounded bg-primary hover:bg-primary-dark"
          on:click|stopPropagation={(e) => (newAddressModal = true) && reset()}><span
          class="material-icons">add</span>
          Add
        </button>
      </div>
      <p>
        By default, mails are sent from a pseudo-randomly generated default
        addresses, like
        <strong><i>vldsbgfr+user@mailwhale.dev</i></strong>
        or so. Alternatively, you can specify custom sender addresses from a
        domain that you own. However, there are a few conditions to be met to do
        so.
      </p>
      <br/>
      <p>
        First, you have to
        <strong>verify</strong>
        the respective e-mail address. That is, once you specify it below, a
        confirmation mail is sent to that address. Only after successful
        verification it can be used to send mail.
      </p>
      <br/>
      <p>
        Second, you have to set proper
        <a
          href="https://blog.mailtrap.io/spf-records-explained/"
          target="_blank"
          class="text-primary">SPF</a>
        and
        <a
          href="https://blog.mailtrap.io/dkim/"
          target="_blank"
          class="text-primary">DKIM</a>
        records for your domain, which permit MailWhale to send mail on your
        behalf. SPF and DKIM are security measures in the context of e-mail that
        aim at verifying the real sender of a message and prevent spam. In order
        to send mail with a custom domain as part of the sender address, you
        need to provide certain SPF- and DKIM DNS records for that domain, which
        will subsequently be verified by the recipient mail server.
      </p>
      <br/>
      <p>
        For
        <strong>SPF</strong>
        please refer to
        <a
          href="https://github.com/muety/mailwhale#spf-check"
          target="_blank"
          class="text-primary">this README section</a>
        on GitHub.
        <strong>DKIM</strong>
        is not implemented, yet.
      </p>
      <br>
      <p>
        Chances are, depending on the recipient's mail provider, that you will also have to set proper <strong><a href="https://blog.mailtrap.io/dmarc-explained/" target="_blank">DMARC</a></strong> records in order for your mails to not be considered spam.
      </p>

      <h2 class="mt-8 mb-4 text-lg font-semibold">Your Addresses</h2>
      {#if me && me.senders && me.senders.length}
        <div class="flex flex-col space-y-4">
          {#each me.senders as sender, i}
            <div class="address-card">
              <div class="info">
                <span class="text-sm font-semibold">#{i + 1}</span>
                <span class="p-1" title="Client ID">{sender.mail} </span>
                <div class="badges">
                  {#if sender.verified}
                    <span class="bg-green-600">verified</span>
                  {:else}<span class="bg-red-600">not verified</span>{/if}
                </div>
              </div>
              <div>
                <a
                  class="text-sm underline cursor-pointer text-primary hover:text-primary-dark"
                  on:click={confirm('Are you sure?') && _deleteAddress(sender)}>Remove</a>
              </div>
            </div>
          {/each}
        </div>
      {:else}
        <div
          class="flex items-center justify-center w-full py-12 text-gray-500">
          <i>No addresses added, yet.</i>
        </div>
      {/if}
    </div>

    {#if newAddressModal}
      <Modal on:close={(e) => (newAddressModal = false) || reset()}>
        <h1 class="text-2xl font-semibold" slot="header">
          Add new sender address
        </h1>
        <div slot="main" style="min-width: 400px;">
          <form
            class="flex flex-col w-full space-y-4"
            on:submit|preventDefault={_addAddress}>
            <div class="flex flex-col w-full space-y-1">
              <label for="mail-input" class="font-semibold">E-Mail</label>
              <input
                type="text"
                class="p-2 border-2 rounded-md border-primary"
                name="mail-input"
                placeholder="E.g. My App Admin <admin@example.org>"
                required
                bind:value={newAddress}/>
            </div>

            <div class="flex justify-between py-2">
              <div/>
              <button
                type="submit"
                class="px-4 py-2 text-white rounded-md bg-primary hover:bg-primary-dark">Add
                & Verify
              </button>
            </div>
          </form>
        </div>
      </Modal>
    {/if}
  </div>
</Layout>
