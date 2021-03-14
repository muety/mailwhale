<script>
  import Layout from '../layouts/Main.svelte'
  import Navigation from '../components/Navigation.svelte'

  import { successes } from '../stores/alerts'
  import { sendMail } from '../api/mails'

  const emptyMail = {
    to: '',
    from: '',
    subject: '',
    body: '',
    html: false,
  }

  let newMail

  reset()

  async function _sendMail() {
    await sendMail({
      from: newMail.from,
      to: [newMail.to],
      subject: newMail.subject,
      html: newMail.html ? newMail.body : null,
      text: !newMail.html ? newMail.body : null,
    })
    successes.spawn('Mail sent successfully!')
    reset()
  }

  function reset() {
    newMail = JSON.parse(JSON.stringify(emptyMail))
  }
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
          <label for="to-input" class="w-1/4">Recipient:</label>
          <input
            type="to"
            class="border-2 border-primary rounded-md p-2 flex-grow"
            name="to-input"
            placeholder="E.g. 'John Doe <john@example.org>'"
            required
            bind:value={newMail.to} />
        </div>

        <div class="flex space-x-2 items-center">
          <label for="from-input" class="w-1/4">Sender:</label>
          <input
            type="from"
            class="border-2 border-primary rounded-md p-2 flex-grow"
            name="from-input"
            placeholder="E.g. 'Jane Doe <jane@example.org>'"
            required
            bind:value={newMail.from} />
        </div>

        <div class="flex space-x-2 items-center">
          <label for="subject-input" class="w-1/4">Subject:</label>
          <input
            type="subject"
            class="border-2 border-primary rounded-md p-2 flex-grow"
            name="subject-input"
            placeholder="Some subject line"
            required
            bind:value={newMail.subject} />
        </div>

        <div class="flex flex-col mt-4 space-y-2">
          <label for="message-input">Text</label>
          <textarea
            class="border-2 border-primary rounded-md p-2 flex-grow w-full"
            style="min-height: 300px;"
            placeholder="Your message"
            required
            name="message-input"
            bind:value={newMail.body} />
        </div>

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
              class="material-icons">send</span>
            Send</button>
        </div>
      </form>
    </div>
  </div>
</Layout>
