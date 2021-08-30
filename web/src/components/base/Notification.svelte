<script lang="typescript">
	import { onMount } from 'svelte';

	import type { NotificationType } from '../../service/types';

	import Icon from './Icon.svelte';

	export let message: string = '';
	export let type: NotificationType = 'info';
	export let duration: number;

	let node: HTMLElement;

	onMount(() => {
		if (duration) {
			setTimeout(() => {
				node.remove();
			}, duration);
		}
	});
</script>

<div
	class="w-full text-white rounded shadow p-2 flex items-center space-x-2"
	class:bg-danger-200={type === 'danger'}
	class:bg-info={type === 'info'}
	class:bg-warning={type === 'warning'}
	class:bg-success={type === 'success'}
	class:text-yellow-700={type === 'warning'}
	class:text-lime-800={type === 'success'}
	bind:this={node}
>
	<div class="w-6">
		{#if type === 'danger'}
			<Icon icon="alert-octagon" />
		{:else if type === 'info'}
			<Icon icon="info" />
		{:else if type === 'warning'}
			<Icon icon="alert-triangle" />
		{:else if type === 'success'}
			<Icon icon="check-circle" />
		{/if}
	</div>
	<div>
		<h2 class="text-lg">
			{message}
		</h2>
	</div>
</div>
