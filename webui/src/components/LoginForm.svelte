<script>
    import { createEventDispatcher } from "svelte";

    const dispatch = createEventDispatcher()

    let username, password;

    function login() {
        if (!username || !password) return
        dispatch('login', { username, password })
    }
</script>

<form class="p-4 w-full flex flex-col space-y-4" on:submit|preventDefault="{login}">
    <div class="flex flex-col w-full space-y-1">
        <label for="email-input">E-Mail</label>
        <input
            type="email"
            class="border-2 border-primary rounded-md p-2"
            name="email-input"
            placeholder="john.doe@example.org"
            required
            bind:value={username} />
    </div>

    <div class="flex flex-col w-full space-y-1">
        <label for="password-input">Password</label>
        <input
            type="password"
            class="border-2 border-primary rounded-md p-2"
            name="password-input"
            placeholder="********"
            required
            bind:value={password} />
    </div>

    <div class="flex justify-between py-2">
        <div />
        <button
            type="submit"
            class="py-2 px-4 text-white bg-primary rounded-md hover:bg-primary-dark">Log In</button>
    </div>
</form>
