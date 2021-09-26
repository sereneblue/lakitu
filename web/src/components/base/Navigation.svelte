<script lang="ts">
	import { slide } from 'svelte/transition';
	import { page } from '$app/stores';

	import Icon from './Icon.svelte';

	let icon: string;
	let show: boolean = false;

	const path: string = $page.path;

	const handleMenuClick = (): void => {
		show = !show;
	};

	$: icon = show ? 'x' : 'cloud';
</script>

<div class="relative z-20">
	{#if show}
		<div
			transition:slide
			class="fixed right-8 bottom-16 mb-2 rounded-md shadow bg-white overflow-hidden focus:outline-none"
			role="menu"
			aria-orientation="vertical"
			aria-labelledby="menu-button"
			tabindex="-1"
		>
			<div role="none">
				<a
					href="machines"
					class="text-gray-700 px-4 py-2 hover:bg-nord6 block py-1 text-sm no-underline flex items-center space-x-2"
					class:hidden={path === '/machines'}
					role="menuitem"
					tabindex="-1"
					rel="external"
				>
					<span class="w-4 h-4">
						<Icon icon="monitor" />
					</span>
					<span>Machines</span>
				</a>
				<a
					href="settings"
					class="text-gray-700 px-4 py-2 hover:bg-nord6 block py-1 text-sm no-underline flex items-center space-x-2"
					class:hidden={path === '/settings'}
					role="menuitem"
					tabindex="-1"
					rel="external"
				>
					<span class="w-4 h-4">
						<Icon icon="settings" />
					</span>
					<span>Settings</span>
				</a>
				<a
					href="logout"
					class="text-gray-700 px-4 py-2 hover:bg-nord6 block py-1 text-sm no-underline flex items-center space-x-2"
					role="menuitem"
					tabindex="-1"
					rel="external"
				>
					<span class="w-4 h-4">
						<Icon icon="log-out" />
					</span>
					<span>Logout</span>
				</a>
			</div>
		</div>
	{/if}

	<button
		on:click={handleMenuClick}
		class="fixed p-2 right-6 bottom-6 mr-2 rounded-md text-white dark:text-nord4 bg-accent-100 hover:bg-accent-200 dark:bg-accent-200 dark:hover:bg-accent-300"
	>
		<div class="w-6 h-6 cursor-pointer">
			<Icon {icon} />
		</div>
	</button>
</div>
