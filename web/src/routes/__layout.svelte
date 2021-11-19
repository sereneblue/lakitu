<script context="module" lang="ts">
	import type { LoadInput, LoadOutput } from '@sveltejs/kit/types.internal';

	export async function load({ page, fetch }: LoadInput): Promise<LoadOutput> {
		let res = await fetch('http://' + page.host + '/session/loggedin');
		let data = await res.json();

		const loggedIn = data.success;

		if (loggedIn) {
			res = await fetch('http://' + page.host + '/jobs');
			let pendingData = await res.json();
	
			const hasPending = pendingData.data != undefined && pendingData.data.isComplete == false;

			if (hasPending && page.path != '/loading') {
				return { status: 302, redirect: '/loading' };
			} 
		} else if (!loggedIn && page.path != '/login' && page.path != '/setup' && page.path != '/') {
			return { status: 302, redirect: '/login' };
		}

		return {
			status: 200
		}
	}
</script>

<div>
	<div class="flex bg-nord0 text-nord4 h-full min-h-screen px-6 py-8">
		<slot />
	</div>
</div>

<style lang="postcss" global>
	@tailwind base;
	@tailwind components;
	@tailwind utilities;

	a {
		@apply cursor-pointer
			   underline
			   text-accent-200 hover:text-accent-300 font-semibold;
	}
</style>
