<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';

	import { Button, Icon, Navigation } from '../components/base';
	import { Machine } from '../components/machines';
	import { 
		CreateMachine, 
		DeleteMachine, 
		MachineDetails, 
		ResizeMachine, 
		StartMachine, 
		StopMachine, 
		TransferMachine 
	} from '../components/modals';

	import type { MachineEventData } from '../service/types';

	import { getRegions, getMachines } from '../service/machines';

	let state = {
		currentMachine: null,
		loading: false,
		machines: [],
		regions: [],
		modal: {
			showCreate: false,
			showDelete: false,
			showDetails: false,
			showResize: false,
			showStart: false,
			showStop: false,
			showTransfer: false,
		}
	};

	const toggleModal = (modalName: string, isOpen: boolean): void => {
		state.modal[modalName] = isOpen;
	}

	const updateMachineList = async (): Promise<void> => {
		state.loading = true;

		let res = await getMachines();

		if (res.success) {
			state.machines = res.data.machines;
		}

		state.loading = false;
	}

	const handleMachineEvent = (e: Event): void => {
		state.currentMachine = (e.detail as MachineEventData).machine;
		switch ((e.detail as MachineEventData).event) {
			case 'delete':
				toggleModal('showDelete', true);
				break;
			case 'details':
				toggleModal('showDetails', true);
				break;
			case 'resize':
				toggleModal('showResize', true);
				break;
			case 'start':
				toggleModal('showStart', true);
				break;
			case 'stop':
				toggleModal('showStop', true);
				break;
			case 'transfer':
				toggleModal('showTransfer', true);
				break;
			default:
				break;
		}
	}

	onMount(async () => {
		let regions = [];
		let machines = [];

		state.loading = true;
		
		let res = await Promise.all([getRegions(), updateMachineList()]);

		if (res[0].success) {
			state.regions = res[0].data.regions.map((r) => {
				return {
					text: r.name,
					value: r.id
				};
			});
		}

		setInterval(updateMachineList, 60 * 1000);
	})
</script>

<svelte:head>
	<title>Machines</title>
</svelte:head>

<Navigation />

<CreateMachine 
	show={state.modal.showCreate} 
	regions={state.regions}
	on:close={() => toggleModal('showCreate', false)}
	on:refresh={updateMachineList} />
<DeleteMachine
	show={state.modal.showDelete}
	machine={state.currentMachine}
	on:close={() => toggleModal('showDelete', false)}
	on:refresh={updateMachineList} />
<MachineDetails
	show={state.modal.showDetails}
	machine={state.currentMachine}
	on:close={() => toggleModal('showDetails', false)} />
<ResizeMachine
	show={state.modal.showResize}
	machine={state.currentMachine}
	on:close={() => toggleModal('showResize', false)}
	on:refresh={updateMachineList} />
<StartMachine
	show={state.modal.showStart}
	machine={state.currentMachine}
	on:close={() => toggleModal('showStart', false)}
	on:refresh={updateMachineList} />
<StopMachine
	show={state.modal.showStop}
	machine={state.currentMachine}
	on:close={() => toggleModal('showStop', false)}
	on:refresh={updateMachineList} />
<TransferMachine
	show={state.modal.showTransfer}
	machine={state.currentMachine}
	regions={state.regions}
	on:close={() => toggleModal('showTransfer', false)}
	on:refresh={updateMachineList} />

<div class="flex flex-col items-center w-full">
	<div class="min-w-full max-w-screen-lg h-full my-4">
		<div class="max-w-screen-xl mx-auto">
			<div class="mt-8">
				<div class="flex justify-between items-center mb-4">
					<h2 class="flex space-x-2 items-center text-2xl mb-4">
						<div>
							Machines
						</div>
						{#if state.loading }
							<div class="w-5 h-5">
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
					</h2>
					{#if state.machines.length}
						<Button icon="plus" text="Create" onClick={() => toggleModal('showCreate', true)} />
					{/if}
				</div>
				<div>
					{#if state.machines.length}
						<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
							{#each state.machines as m}
								<Machine 
									data={m}
									on:emit={handleMachineEvent} />
							{/each}
						</div>
					{:else}
						<div class="text-center w-full">
							<div class="w-full flex justify-center">
								<div class="w-52 h-52 opacity-25">
									<Icon icon="server" stroke={1} />
								</div>
							</div>
							<div class="my-4 text-lg">
								<div class="font-bold">No available machines</div>
								<div class="text-sm opacity-75">Create a machine below</div>
							</div>
							<div>
								<Button icon="plus" text="Create" onClick={() => toggleModal('showCreate', true)} />
							</div>
						</div>
					{/if}
				</div>
			</div>
		</div>
	</div>
</div>
