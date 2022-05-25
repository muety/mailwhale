<script>
  import LoginForm from '../components/LoginForm.svelte'
  import Layout from '../layouts/Main.svelte'

  import { user } from '../stores/auth'

  async function login({ detail }) {
    try {
        await user.login(detail)
        window.location.replace('clients')
    } catch (e) {}
  }
</script>

<Layout>
  <div slot="content" class="flex flex-col items-center justify-center">
    <h1 class="text-2xl font-semibold mb-4">Login</h1>
    <p class="w-1/2 lg:w-1/3 text-center">If you have problems logging in, make sure you have confirmed the activation mail (in case the server requires user verification).</p>
    <div class="w-1/2 py-8 lg:w-1/3">
      <LoginForm on:login={login} />
    </div>
  </div>
</Layout>
