<script>
  import LoginForm from '../components/LoginForm.svelte'
  import Layout from '../layouts/Main.svelte'

  import { user } from '../stores/auth'

  function login({ detail }) {
    user.login(detail)
    window.location.replace('clients')
  }
</script>

<style global lang="postcss">
  :root {
    --color-primary: #159ce4;
    --color-primary-dark: #138dce;
    --color-text: #4b5563;
  }

  @tailwind base;
  @tailwind components;
  @tailwind utilities;
</style>

<Layout>
  <div slot="content" class="flex flex-col items-center justify-center">
    <h1 class="text-2xl font-semibold">Login</h1>
    <div class="py-8 w-1/2 lg:w-1/3">
      <LoginForm on:login={login} />
    </div>
  </div>
</Layout>
