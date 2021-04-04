<script>
  import Layout from '../layouts/Main.svelte'
  import Navigation from '../components/Navigation.svelte'
</script>

<Layout>
<div slot="content" class="flex">
  <div class="w-1/4">
    <Navigation />
  </div>
  <div class="max-w-screen-lg">
    <h1 class="mb-4 text-2xl font-semibold">SPF & DKIM Configuration</h1>
    <a href="https://blog.mailtrap.io/spf-records-explained/" target="_blank" class="text-primary">SPF</a> and <a href="https://blog.mailtrap.io/dkim/" target="_blank" class="text-primary">DKIM</a> are security measures in the context of e-mail that aim at verifying the real sender of a message and prevent spam. In order to send mail with a custom domain as part of the sender address, you need to provide certain SPF- and DKIM DNS records for that domain, which will subsequently be verified by the recipient mail server.

    <h2 class="mt-8 text-lg font-semibold">SPF</h2>
    <p>
      Please refer to <a href="https://github.com/muety/mailwhale#spf-check" target="_blank" class="text-primary">this README section</a> on GitHub.
    </p>

    <h2 class="mt-8 text-lg font-semibold">DKIM</h2>
    <p>
      DKIM verification is not yet implemented.
    </p>
  </div>
</div>

</Layout>
