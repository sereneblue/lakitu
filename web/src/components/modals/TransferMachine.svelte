<script lang="ts">
	import { createEventDispatcher } from 'svelte'; 
	import { Modal } from '../base';
	import { SelectField } from '../fields';

	import type { MachineData, ModalAction } from '../../service/types';
	import { transferMachine } from '../../service/machines';
	import { showJobProgress } from '../../service/util';

	export let show: boolean = false;
	export let machine: MachineData;
	export let regions: string[] = [];

	let dismissible = true;
	let jobStatusEl: HTMLElement;
	
	let state = {
		region: '',
		regions: []
	}

	const dispatch = createEventDispatcher();

	const closeModal = () => {
		show = false;
		dispatch('close');
	}

	$: {
		if (machine) {
			state.regions = regions.filter(r => r.value != machine.region)
		}
	};

	let action: ModalAction = {
		text: 'Transfer',
		isDisabled: true,
		func: async () => {
			dismissible = false;

			let res = await transferMachine(machine.uuid, state.region);

			if (res.success) {
			   return new Promise((resolve) => {
					showJobProgress(jobStatusEl, res.data.jobId, () => {
						dismissible = true;
						setTimeout(() => {
							dispatch('refresh');
							closeModal();
							resolve();
						}, 3000);
					})
			   });
			}
		}
	}

	$: action.isDisabled = state.region == '';  
</script>

<Modal title="Change region" type="warning" {action} {dismissible} {show} on:close={closeModal}>	
	<div class="mt-2 flex flex-col space-y-1">
		<div class="flex flex-col">
			<div class="font-bold text-sm tracking-wide opacity-50">Transfer from</div>
			<div class="text-xl -mt-1">{machine.region}</div>
		</div>
		<div class="flex flex-col pt-2">
			<div class="font-bold text-sm tracking-wide opacity-50">Transfer to</div>
			<SelectField
				on:change={(e) => (state.region = e.detail.value)}
				placeholder="AWS Region"
				selectedValue={''}
				options={state.regions}
				full
			/>
		</div>
		<div class="pt-4">
			<div bind:this={jobStatusEl} />
		</div>
	</div>
</Modal>