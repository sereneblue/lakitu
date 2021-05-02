<script lang="typescript">
	import Icon from './Icon.svelte';

	type ButtonSize = 'xs' | 'sm' | 'md' | 'lg';
	type ButtonType = 'default' | 'danger' | 'primary';

	export let text: string = '';
	export let icon: string = '';
	export let type: ButtonType = 'default';
	export let size: ButtonSize = 'md';
	export let disabled: boolean = false;
	export let loading: boolean = false;
	export let outline: boolean = false;
	export let full: boolean = false;

	export let onClick: Function = () => {};

	const btnClick = (): void => {
		if (!loading) {
			onClick();
		}
	};
</script>

<button
	class="
		text-center text-white font-bold 
		py-1 px-2
		relative
		rounded shadow
		focus:outline-none focus-visible:ring-2
		transition transition-colors"
	class:xs={size === 'xs'}
	class:sm={size === 'sm'}
	class:md={size === 'md'}
	class:lg={size === 'lg'}
	class:disabled
	class:outline
	class:w-full={full}
	class:default={type === 'default'}
	class:danger={type === 'danger'}
	class:primary={type === 'primary'}
	on:click={btnClick}
	{disabled}
>
	<div class="flex items-center w-full justify-center" class:invisible={loading}>
		{#if icon}
			<span class="icon">
				<Icon {icon} />
			</span>
		{/if}
		{#if text}
			<span>{text}</span>
		{/if}
	</div>
	{#if loading}
		<!--- loading spinner --->
		<div class="absolute top-0 left-0 flex flex-col items-center justify-center w-full h-full">
			<span class="w-full h-3/5">
				<svg
					class="animate-spin w-full h-full text-white"
					xmlns="http://www.w3.org/2000/svg"
					fill="none"
					viewBox="0 0 24 24"
				>
					<circle
						class="opacity-25"
						cx="12"
						cy="12"
						r="10"
						stroke="currentColor"
						stroke-width="4"
					/>
					<path
						class="opacity-75"
						fill="currentColor"
						d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
					/>
				</svg>
			</span>
		</div>
	{/if}
</button>

<style type="postcss">
	/* size */

	.xs {
		@apply text-xs;
	}
	.sm {
		@apply text-sm;
	}
	.md {
		@apply text-base;
	}
	.lg {
		@apply text-lg px-4;
	}

	/* color */

	.disabled {
		@apply opacity-50 cursor-not-allowed;
	}

	.default {
		@apply bg-nord3 hover:bg-nord2 dark:bg-nord1 dark:hover:bg-nord2
			ring-accent-100 dark:ring-accent-200 ring-opacity-50;
	}
	.default.disabled {
		@apply hover:bg-nord3 dark:hover:bg-nord1;
	}
	.default.outline {
		@apply bg-transparent hover:bg-nord3 dark:hover:bg-nord1
			text-nord3 dark:text-nord1 hover:text-white
			border border-nord3 dark:border-nord1;
	}
	.default.outline.disabled {
		@apply hover:bg-transparent hover:text-nord3 dark:hover:text-nord1;
	}
	.primary {
		@apply bg-accent-100 hover:bg-accent-200 dark:bg-accent-200 dark:hover:bg-accent-300
			   ring-accent-100 dark:ring-accent-200 ring-opacity-50;
	}
	.primary.disabled {
		@apply hover:bg-accent-100 dark:hover:bg-accent-200;
	}
	.primary.outline {
		@apply bg-transparent hover:bg-accent-100 dark:hover:bg-accent-200 
			text-accent-100 dark:text-accent-200 hover:text-white
			border border-accent-100 dark:border-accent-200;
	}
	.primary.outline.disabled {
		@apply hover:bg-transparent hover:text-accent-100 dark:hover:text-accent-200;
	}

	.danger {
		@apply bg-danger-100 hover:bg-danger-200 dark:bg-danger-200 dark:hover:bg-danger-300
			   ring-danger-100 dark:ring-danger-200 ring-opacity-50;
	}
	.danger.disabled {
		@apply hover:bg-danger-100 dark:hover:bg-danger-200;
	}
	.danger.outline {
		@apply bg-transparent hover:bg-danger-100 dark:hover:bg-danger-200
			text-danger-100 dark:text-danger-200 hover:text-white
			border border-danger-100 dark:border-danger-200;
	}
	.danger.outline.disabled {
		@apply hover:bg-transparent hover:text-danger-100 dark:hover:text-danger-200;
	}

	/* icon */

	span.icon {
		@apply w-4 h-4;
	}

	span.icon + span {
		@apply ml-2;
	}
</style>
