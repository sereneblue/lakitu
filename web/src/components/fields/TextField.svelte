<script lang="ts">
	import feather from 'feather-icons';
	import { createEventDispatcher } from 'svelte';

	import { Icon } from '../base';

	type Align = 'left' | 'center';
	type InputType = 'text' | 'number' | 'password';

	export let icon: string = '';
	export let type: InputType = 'text';
	export let label: string = '';
	export let value: string = '';
	export let error: string = '';
	export let full: boolean = false;
	export let disabled: boolean = false;
	export let align: Align = 'left';
	export let placeholder: string = '';
	export let submit: Function;
	export let suffix: string = '';
	export let stroke: number = 2;

	let loading = false;

	const dispatch = createEventDispatcher();

	const change = (e: Event): void => {
		dispatch('change', {
			value: e.target.value
		});
	};

	const handleEnter = async (e: Event): Promise<void> => {
		if (e.key === 'Enter' && submit && !loading) {
			loading = true;
			await submit();
			loading = false;
		}
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
	{#if icon}
		<div class="absolute ml-2 mt-2 w-6 h-6 opacity-50">
			<Icon {icon} {stroke} />
		</div>
	{/if}
	<div class="flex items-center">
		<input
			class="
				appearance-none
				bg-nord6 bg-opacity-10 
				border-2 border-transparent
				w-full
				p-2
				leading-tight
				focus:outline-none focus:border-accent-200"
			class:w-full={full}
			class:pl-9={icon}
			class:text-left={align === 'left'}
			class:text-center={align === 'center'}
			class:error={hasError}
			class:rounded={!suffix}
			on:input={change}
			on:keypress={handleEnter}
			class:rounded-l={suffix}
			{value}
			{type}
			{placeholder}
			{disabled}
		/>
		{#if suffix}
			<span
				class="
				bg-nord6 bg-opacity-10
				border
				border-l-2
				rounded-r
				p-2
				leading-tight"
			>
				<span class="opacity-75">
					{suffix}
				</span>
			</span>
		{/if}
	</div>
	<div
		class="text-danger-200"
		class:text-left={align === 'left'}
		class:text-right={align === 'right'}
		class:text-center={align === 'center'}
	>
		{error}
	</div>
</label>

<style type="postcss">
	.error {
		@apply border-danger-200;
	}

	.modal-content label input {
		@apply bg-nord0 bg-opacity-5;
	}
</style>
