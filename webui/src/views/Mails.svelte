<script>
  import { onMount } from 'svelte'

  import Layout from '../layouts/Main.svelte'
  import Navigation from '../components/Navigation.svelte'

  import { successes } from '../stores/alerts'
  import { sendMail } from '../api/mails'
  import { getTemplates } from '../api/template'
  import { getClients } from '../api/clients'

  import { extractVars } from '../utils/template'

  const emptyMail = {
    to: '',
    subject: '',
    body: '',
    html: false,
    template_id: null,
    template_vars: null,
  }

  let newMail

  let templates = []
  let selectedTemplate
  let templateVarsStr = '{}'

  let clients = []
  let selectedClient

  let sending = false

  reset()

  async function _sendMail() {
    sending = true
    try {
      await sendMail({
        to: [newMail.to],
        subject: newMail.subject,
        html: newMail.html && newMail.body ? newMail.body : null,
        text: !newMail.html && newMail.body ? newMail.body : null,
        template_id: selectedTemplate ? selectedTemplate.id : null,
        template_vars: selectedTemplate
          ? JSON.parse(templateVarsStr || '{}')
          : null,
      }, selectedClient?.id)
      successes.spawn('Mail sent successfully!')
    } finally {
      reset()
    }
  }

  function onTemplateSelected() {
    if (!selectedTemplate) return
    templateVarsStr = JSON.stringify(
      extractVars(selectedTemplate.content),
      null,
      4
    )
  }

  function reset() {
    newMail = JSON.parse(JSON.stringify(emptyMail))
    sending = false
  }

  onMount(() => {
    getTemplates()
      .then((result) => (templates = result))
      .catch(() => {})

    getClients()
      .then((result) => (clients = result.filter(c => c.permissions.includes('send_mail'))))
      .catch(() => {})
  })
</script>

<Layout>
  <div slot="content" class="flex">
    <div class="w-1/4">
      <Navigation />
    </div>
    <div class="flex flex-col px-12 w-full w-3/4">
      <div class="flex justify-between mb-8">
        <h1 class="text-2xl font-semibold">Send Test E-Mail</h1>
      </div>

      <form
        class="w-full flex flex-col space-y-4"
        on:submit|preventDefault={_sendMail}>
        <div class="flex space-x-2 items-center">
          <label for="to-input" class="w-1/4 font-semibold">Recipient:</label>
          <input
            type="text"
            class="border-2 border-primary rounded-md p-2 flex-grow"
            name="to-input"
            placeholder="E.g. 'John Doe <john@example.org>'"
            required
            bind:value={newMail.to} />
        </div>

        <div class="flex space-x-2 items-center">
          <label
            for="subject-input"
            class="w-1/4 font-semibold">Subject:</label>
          <input
            type="text"
            class="border-2 border-primary rounded-md p-2 flex-grow"
            name="subject-input"
            placeholder="Some subject line"
            required
            bind:value={newMail.subject} />
        </div>

        <div class="flex space-x-2 items-center">
          <label
            for="template-input"
            class="w-1/4 font-semibold">Template:</label>
          <select
            class="border-2 border-primary rounded-md p-2 flex-grow cursor-pointer"
            bind:value={selectedTemplate}
            on:change={onTemplateSelected}>
            <option selected value>No template</option>
            {#each templates as template}
              <option value={template}>{template.name} ({template.id})</option>
            {/each}
          </select>
        </div>

        <div class="flex space-x-2 items-center">
          <label
            for="template-input"
            class="w-1/4 font-semibold">Send as API Client:</label>
          <select
            class="border-2 border-primary rounded-md p-2 flex-grow cursor-pointer"
            bind:value={selectedClient}>
            <option selected value>None (send as user)</option>
            {#each clients as client}
              <option value={client}>{client.id} ({client.description})</option>
            {/each}
          </select>
        </div>

        {#if !selectedTemplate}
          <div class="flex flex-col mt-4 space-y-2">
            <label for="message-input" class="font-semibold">Text</label>
            <textarea
              class="border-2 border-primary rounded-md p-2 flex-grow w-full"
              style="min-height: 300px;"
              placeholder="Your message (text or HTML)"
              name="message-input"
              bind:value={newMail.body} />
          </div>
        {:else}
          <div class="flex flex-col mt-4 space-y-2">
            <label for="vars-input" class="font-semibold">Template Variables
              (JSON)</label>
            <textarea
              class="border-2 border-primary rounded-md p-2 flex-grow w-full"
              style="min-height: 300px;"
              placeholder="JSON object of variables to be used in the template"
              name="vars-input"
              bind:value={templateVarsStr} />
          </div>
        {/if}

        <div class="flex justify-between py-2">
          <div class="flex space-x-2 items-center">
            <input
              type="checkbox"
              name="html-checkbox"
              bind:checked={newMail.html} />
            <label for="html-checkbox">HTML ?</label>
          </div>
          <button
            type="submit"
            class="flex items-center px-4 py-2 bg-primary text-white rounded hover:bg-primary-dark"><span
              class="material-icons"
              disabled={sending}>send</span>
            {sending ? 'Sending ...' : 'Send'}</button>
        </div>
      </form>
    </div>
  </div>
</Layout>
