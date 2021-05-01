<script lang="typescript">
	import { createEventDispatcher } from 'svelte';

	type SelectOption = {
		text: string;
		value: number | string;
	};

	export let label: string = '';
	export let selectedValue: string = '';
	export let full: boolean = false;
	export let options: SelectOption[] = [];
	export let placeholder: string = '';

	const dispatch = createEventDispatcher();

	const change = (e: Event): void => {
		selectedValue = e.target.value;
		dispatch('change', {
			value: e.target.value
		});
	};
</script>

<label>
	{#if label}
		<div>{label}</div>
	{/if}
	<select
		class="
			block
			appearance-none
			bg-nord0 bg-opacity-5 dark:bg-nord6 dark:bg-opacity-10 
			border-2 border-transparent
			rounded w-full
			py-1 px-2
			leading-tight
			focus-visible:bg-white
			focus:outline-none focus:border-accent-200"
		class:default={placeholder && selectedValue === ''}
		class:w-full={full}
		on:input={change}
		value={selectedValue}
	>
		{#if placeholder}
			<option value="" disabled selected={selectedValue === ''}>{placeholder}</option>
		{/if}
		{#each options as o}
			<option value={o.value} selected={o.value == selectedValue}>{o.text}</option>
		{/each}
	</select>
</label>

<style type="postcss">
	select {
		background-image: url("data:image/svg+xml,%3csvg xmlns='http://www.w3.org/2000/svg' fill='none' viewBox='0 0 20 20'%3e%3cpath stroke='%236b7280' stroke-linecap='round' stroke-linejoin='round' stroke-width='1.5' d='M6 8l4 4 4-4'/%3e%3c/svg%3e");
		background-position: right 0.5rem center;
		background-repeat: no-repeat;
		background-size: 1.5em 1.5em;
		padding-right: 2.5rem;
		color-adjust: exact;
	}

	.default {
		@apply text-gray-400 dark:text-gray-500;
	}
</style>
