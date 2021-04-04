<script>
    import { createEventDispatcher } from 'svelte';
    import { errors } from '../stores/alerts'
    
    const dispatch = createEventDispatcher()

    let username, password, passwordRepeat;

    function signup() {
        if (!username || !password || !passwordRepeat) return

        if (password !== passwordRepeat) {
            return errors.spawn('Passwords do not match')
        }

        dispatch('signup', { username, password })
    }
</script>

<form class="flex flex-col w-full p-4 space-y-4" on:submit|preventDefault="{signup}">
    <div class="flex flex-col w-full space-y-1">
        <label for="email-input">E-Mail</label>
        <input
            type="email"
            class="p-2 border-2 rounded-md border-primary"
            name="email-input"
            placeholder="john.doe@example.org"
            required
            bind:value={username} />
    </div>

    <div class="flex flex-col w-full space-y-1">
        <label for="password-input">Password</label>
        <input
            type="password"
            class="p-2 border-2 rounded-md border-primary"
            name="password-input"
            placeholder="********"
            required
            bind:value={password} />
    </div>

    <div class="flex flex-col w-full space-y-1">
        <label for="password-input">Password (repeat)</label>
        <input
            type="password"
            class="p-2 border-2 rounded-md border-primary"
            name="password-repeat-input"
            placeholder="********"
            required
            bind:value={passwordRepeat} />
    </div>

    <div class="flex justify-between py-2">
        <div />
        <button
            type="submit"
            class="px-4 py-2 text-white rounded-md bg-primary hover:bg-primary-dark">Sign Up</button>
    </div>
</form>
