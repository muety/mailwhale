<script>
  import SignupForm from '../components/SignupForm.svelte'
  import Layout from '../layouts/Main.svelte'

  import { createUser } from '../api/users'
  import { successes } from '../stores/alerts'

  async function signup({ detail }) {
    const user = await createUser({ email: detail.username, password: detail.password })
    successes.spawn('Successfully signed up')
    window.location.replace('login')
  }
</script>

<Layout>
  <div slot="content" class="flex flex-col items-center justify-center">
    <h1 class="text-2xl font-semibold">Sign Up</h1>
    <div class="w-1/2 py-8 lg:w-1/3">
      <SignupForm on:signup={signup}/>
    </div>
  </div>
</Layout>
