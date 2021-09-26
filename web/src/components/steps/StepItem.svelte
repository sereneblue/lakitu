<script lang="ts">
	import { createEventDispatcher, getContext } from 'svelte';

	const dispatch = createEventDispatcher();
	const steps = getContext('steps');
	const currentStep = getContext('currentStep');

	const stepIndex = (() => {
		return steps.index();
	})();

	export let label: string = '';
</script>

<div data-step-index={stepIndex} data-label={label}>
	{#if $currentStep === stepIndex}
		<slot />
	{/if}
</div>
