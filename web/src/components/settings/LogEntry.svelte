<script lang="ts">
	import { slide } from 'svelte/transition';

	import { Button, Icon } from '../base';
	import { Status } from '../../service/types';
	import type { Log } from '../../service/types'; 

	export let details: Log;

	let state = {
		dateFmt: new Date(),
		open: false
	}

	const fmtTimestamp = (timestamp: number): string => {
		state.dateFmt.setTime(timestamp);

		return state.dateFmt.toLocaleString(); 
	}
</script>

<div>
	<div class="flex space-x-4 rounded w-full p-2">
		<div class="flex items-center">
			<div class="flex items-center justify-center w-8 h-8">
				{#if details.status == Status.PENDING }
					<div>
						<Icon icon="clock" />
					</div>
				{:else if details.status == Status.COMPLETE}
					<div class="text-accent-200">
						<Icon icon="check-circle" />
					</div>
				{:else if details.status == Status.ERROR}
					<div class="text-danger-200">
						<Icon icon="alert-circle" />
					</div>
				{:else if details.status == Status.CANCELED}
					<div class="text-danger-200">
						<Icon icon="slash" />
					</div>
				{/if}
			</div>
		</div>
		<div class="flex-1">
			<div class="text-2xl">
				{details.event} <span class="text-base opacity-50">(#{details.id})</span>
			</div>
			<div class="opacity-75">
				{fmtTimestamp(details.timestamp)}
			</div>
		</div>
		{#if details.status == Status.ERROR}
			<div class="flex items-center">
				<Button icon={state.open ? "chevron-up" : "chevron-down"} size="md" flat onClick={() => state.open = !state.open} />
			</div>
		{/if}
	</div>
	{#if state.open}
		<p transition:slide>
			{details.errorInfo}
		</p>
	{/if}
</div>