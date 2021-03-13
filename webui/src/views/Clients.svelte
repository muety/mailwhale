<script>
	import Layout from "../layouts/Main.svelte";
	import Navigation from "../components/Navigation.svelte";
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
		<div class="flex px-12">
			<h1 class="text-2xl font-semibold">Manage Clients</h1>
		</div>
	</div>
</Layout>
