<script>
    let currentUser
</script>

<div class="flex space-x-2 items-center h-full">
    {#if currentUser}
        <span class="text-sm">👤 {currentUser}</span>
        <span>|</span>
        <a
            href="#"
            class="text-sm text-primary hover:text-primary-dark">Logout</a>
    {:else}
        <a href="/login" class="text-primary hover:text-primary-dark font-semibold">Login</a>
    {/if}
</div>
