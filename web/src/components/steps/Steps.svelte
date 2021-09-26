<script lang="ts">
	import { onMount, setContext } from 'svelte';
	import { writable } from 'svelte/store';

	export function nextStep() {
		if ($currentStep + 1 <= numSteps) {
			$currentStep = $currentStep + 1;
		}
	}

	let numSteps = 0;
	let currentStep = writable(0);
	let isReady: boolean = false;
	let steps: any;
	let stepItems: HTMLElement[] = [];

	setContext('steps', {
		index: () => {
			numSteps += 1;
			return numSteps;
		}
	});

	setContext('currentStep', currentStep);

	onMount(() => {
		stepItems = steps.querySelectorAll('div[data-step-index]');

		if (stepItems.length) {
			$currentStep = 1;
		}
	});
</script>

<svelte:window on:nextStep={nextStep} />

<div class="w-full max-w-screen-lg">
	<nav class="w-full mx-auto">
		<ul class="flex justify-center">
			{#each stepItems as step}
				<li
					class="flex-1 text-center font-bold"
					class:active={$currentStep == step.dataset.stepIndex}
					class:complete={$currentStep > step.dataset.stepIndex}
				>
					<div class="flex justify-center relative mx-auto mb-1">
						<div class="step-index--num z-10">
							{step.dataset.stepIndex}
						</div>
						<div class="step-index--rail" />
					</div>
					<div class="hidden md:block">
						{step.dataset.label}
					</div>
				</li>
			{/each}
		</ul>
	</nav>
	<section class="py-8" bind:this={steps}>
		<slot />
	</section>
</div>

<style type="postcss" global>
	.step-index--num {
		@apply bg-gray-300 dark:bg-gray-600
        text-lg text-center dark:text-white font-bold
        w-8 h-8
        border border-2 border-nord6 dark:border-nord1
        flex justify-center items-center
        rounded-full;
	}
	.active .step-index--num {
		@apply bg-nord6 dark:bg-nord1
        border-accent-100 dark:border-accent-200;
	}
	.complete .step-index--num {
		@apply bg-accent-100 dark:bg-accent-200
        text-white
        border-transparent
        transition transition-colors;
	}
	.step-index--rail {
		@apply absolute top-1/2 left-1/2
        w-full h-1
        bg-gray-300 dark:bg-gray-600
        transition transition-colors;
	}
	.complete .step-index--rail {
		@apply bg-accent-100 dark:bg-accent-200;
	}
	ul li:last-child .step-index--rail {
		@apply hidden;
	}
</style>
