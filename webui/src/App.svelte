<script>
    import router from "page";
    import Home from "./views/Home.svelte";
    import Clients from "./views/Clients.svelte";

    let page;

    router("/", () => (page = Home));
    router("/login", () => (page = Home));
    router("/clients", () => (page = Clients));

    router.start();
</script>

<div id="app-container" class="container mx-auto flex flex-col flex-grow">
    <svelte:component this={page} />
</div>
