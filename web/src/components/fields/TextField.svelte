<script lang="typescript">
	import feather from 'feather-icons';
	import { createEventDispatcher } from 'svelte';

	type Align = 'left' | 'center';
	type InputType = 'text' | 'password';

	export let type: InputType = 'text';
	export let label: string = '';
	export let value: string = '';
	export let error: string = '';
	export let full: boolean = false;
	export let disabled: boolean = false;
	export let align: Align = 'left';
	export let placeholder: string = '';
	export let stroke: number = 2;

	const dispatch = createEventDispatcher();

	const change = (e: Event): void => {
		dispatch('change', {
			value: e.target.value
		});
	};

	$: hasError = error;
</script>

<label>
	{#if label}
		<div
			class:text-left={align === 'left'}
			class:text-right={align === 'right'}
			class:text-center={align === 'center'}
		>
			{label}
		</div>
	{/if}
	<input
		class="
			appearance-none
			bg-nord0 bg-opacity-5 dark:bg-nord6 dark:bg-opacity-10 
			border-2 border-transparent
			rounded w-full
			p-2
			leading-tight
			focus-visible:bg-white
			focus:outline-none focus:border-accent-200"
		class:w-full={full}
		class:text-left={align === 'left'}
		class:text-center={align === 'center'}
		class:error={hasError}
		on:input={change}
		{value}
		{type}
		{placeholder}
		{disabled}
	/>
	<div
		class="text-danger-100 dark:text-danger-200"
		class:text-left={align === 'left'}
		class:text-right={align === 'right'}
		class:text-center={align === 'center'}
	>
		{error}
	</div>
</label>

<style type="postcss" global>
	.error {
		@apply border-danger-100 dark:border-danger-200;
	}
</style>
