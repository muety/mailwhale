<script>
	import { onMount } from "svelte";

	import Layout from "../layouts/Main.svelte";
	import Navigation from "../components/Navigation.svelte";
	import { getClients } from "../api/clients";

	let clients = [];

	onMount(async () => {
		clients = await getClients();
	});
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
	<div slot="content" class="flex">
		<Navigation />
		<div class="flex flex-col px-12">
			<h1 class="text-2xl font-semibold">Manage Clients</h1>

			<div>
				{#each clients as client}
					<div>{client.id} ({client.description})</div>
				{/each}
			</div>
		</div>
	</div>
</Layout>
