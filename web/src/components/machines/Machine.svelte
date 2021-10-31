<script lang="ts">
	import { createEventDispatcher } from "svelte";
	
	import { Button, Icon } from '../base';
	import type { MachineData, MachineEvent } from '../../service/types';
	import MachineMenu from './MachineMenu.svelte';

	export let data: MachineData;

	let state;
	
	const dispatch = createEventDispatcher();
	
	const emit = (event: MachineEvent): void => {
		dispatch("emit", {
			event,
			machine: data
		});
	}

	$: state = {...data, menuOpen: false };
</script>

<div class="flex items-center relative rounded shadow hover:shadow-md border border-nord2 hover:border-accent-200">
	<div class="h-24 w-24 p-4 bg-nord1">
		<Icon icon="server" stroke={1} />
	</div>
	<MachineMenu 
		status={state.status} 
		bind:open={state.menuOpen}
		on:emit={(e) => emit(e.detail)} />
	<div class="flex-1 space-y-1 ml-2">
		<div class="font-semibold">{state.name}</div>
		<div class="flex items-center text-white">
			{#if state.status === 'offline'}
				<div class="text-sm bg-gray-500 px-2 rounded">Offline</div>
			{:else if state.status === 'online'}
				<div class="text-sm bg-green-400 px-2 rounded">Online</div>
			{:else if state.status === 'unavailable'}
				<div class="text-sm bg-danger-100 px-2 rounded">Unavailable</div>
			{:else if state.status === 'unknown'}
				<div class="text-sm bg-danger-100 px-2 rounded">Unknown</div>
			{/if}
		</div>
		<div class="flex items-center text-sm opacity-50">
			{state.instanceType} / {state.size} GB / {state.region}
		</div>
	</div>
	<div class="h-full">
		<div class="mr-2 w-5 h-5">
			<Button icon="more-vertical" transparent onClick={() => state.menuOpen = !state.menuOpen} />
		</div>
	</div>
</div>