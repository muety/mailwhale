<script>
  import Layout from '../layouts/Main.svelte'

  let lang = 'curl'
</script>

<Layout>
<div slot="content" class="flex flex-col items-center justify-center space-y-8">
    <img src="https://raw.githubusercontent.com/muety/mailwhale/master/assets/screenshot01.png" alt="Screenshot" width="400px" class="-mt-24"/>
    <h1 class="text-4xl font-semibold text-center">
      MailWhale is a <span class="text-primary">bring-your-own-SMTP</span> mail relay 
    </h1>
    <p class="max-w-screen-lg text-center">
      Or, in other words, it is a web service for sending mails via <strong>REST API</strong>. Think of <a href="https://mailgun.com" target="_blank" class="text-primary">Mailgun</a> or <a href="https://sendgrid.com" target="_blank" class="text-primary">SendGrid</a>, except self-hosted and with less features. Or like <a href="https://cuttlefish.io/" target="_blank" class="text-primary">Cuttlefish</a>, but without having to host your own SMTP server. <strong>Plug any SMTP server</strong>, like <a href="https://mail.google.com" target="_blank" class="text-primary">Google Mail</a>, <a href="https://mailbox.org" target="_blank" class="text-primary">Mailbox.org</a>, etc. and start sending mails from within your application via simple HTTP requests.
    </p>
    <div class="w-full max-w-screen-lg">
      <div class="flex mx-2 mb-1 space-x-4 text-sm font-semibold text-gray-400">
        <span class="cursor-pointer" class:text-gray-500={ lang === 'curl' } on:click={ e => (lang = 'curl') }>cURL</span>
        <span class="cursor-pointer" class:text-gray-500={ lang === 'go' } on:click={ e => (lang = 'go') }>Go</span>
        <span class="cursor-pointer" class:text-gray-500={ lang === 'node' } on:click={ e => (lang = 'node') }>NodeJS</span>
      </div>

      {#if lang === 'curl'}
        <textarea class="w-full p-4 font-mono text-xs text-gray-100 bg-gray-800 rounded-md resize-none" rows="8" disabled>
curl -XPOST -u '<clientId>:<clientSecret>' -H 'content-type: application/json' \
  --data-raw {'{'}'
      "to": ["Jane Doe <jane@doe.com>"],
      "subject": "Dinner tonight?",
      "text": "Hey you! Wanna have dinner tonight?"
  {'}'}' \
  'https://mailwhale.dev/api/mail'
        </textarea>

      {:else if lang === 'go'}
        <textarea class="w-full p-4 font-mono text-xs text-gray-100 bg-gray-800 rounded-md resize-none" rows="22" disabled>
type PlainTextMail struct {'{'}
  To      []string `json:"to"`
  Subject string `json:"subject"`
  Text    string `json:"text"`
{'}'}

mail := PlainTextMail{'{'}
  To:      []string {'{'} "Jane Doe <jane@doe.com>" {'}'},
  Subject: "Dinner tonight?",
  Text:    "Hey you! Wanna have dinner tonight?",
{'}'}

payload, _ := json.Marshal(mail)

req, _ := http.NewRequest("https://mailwhale.dev/api/mail", "application/json", bytes.NewBuffer(payload))
req.SetBasicAuth("<clientId>", "<clientSecret>")

client := &http.Client{'{'}}{'}'}
client.Do(req)
        </textarea>

        {:else if lang === 'node'}
        <textarea class="w-full p-4 font-mono text-xs text-gray-100 bg-gray-800 rounded-md resize-none" rows="16" disabled>
const fetch = require('node-fetch')

const token = Buffer.from('<clientId>:<clientSecret>').toString('base64')

const payload = {'{'}
    to: [ 'Jane Doe <jane@doe.com>' ],
    subject: 'Dinner tonight?',
    text: 'Hey you! Wanna have dinner tonight?'
{'}'}

await fetch('https://mailwhale.dev/api/mail', {'{'}
    method: 'post',
    body: JSON.stringify(payload),
    headers: {'{'} 'Authorization': `Basic ${'{'}token{'}'}` {'}'}
{'}'})
        </textarea>
      {/if}
    </div>

    <div class="flex flex-col items-center pt-8 space-y-4">
      <h2 class="text-xl font-semibold text-primary">Features</h2>
      <ul>
        <li class="flex items-center space-x-1"><span class="material-icons">check</span><span>100 % free and open-source</span></li>
        <li class="flex items-center space-x-1"><span class="material-icons">check</span><span>Built by developers for developers</span></li>
        <li class="flex items-center space-x-1"><span class="material-icons">check</span><span>Self-hosted</span></li>
        <li class="flex items-center space-x-1"><span class="material-icons">check</span><span>Intuitive REST API</span></li>
        <li class="flex items-center space-x-1"><span class="material-icons">check</span><span>Clean web UI</span></li>
        <li class="flex items-center space-x-1"><span class="material-icons">check</span><span>Easy setup</span></li>
        <li class="flex items-center space-x-1"><span class="material-icons">check</span><span>Multi-user support</span></li>
        <li class="flex items-center space-x-1"><span class="material-icons">check</span><span>HTML e-mail templates</span></li>
      </ul>
  </div>
</div>

</Layout>
