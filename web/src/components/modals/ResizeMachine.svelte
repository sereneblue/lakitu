<script lang="ts">
	import { createEventDispatcher } from 'svelte'; 
	import { Modal } from '../base';
	import { TextField } from '../fields';

	import type { MachineData, ModalAction } from '../../service/types';
	import { resizeMachine } from '../../service/machines';
	import { showJobProgress } from '../../service/util';

	export let show: boolean = false;
	export let machine: MachineData;

	let dismissible = true;
	let jobStatusEl: HTMLElement;
	
	let state = {
		size: '',
	}

	const dispatch = createEventDispatcher();

	const closeModal = () => {
		show = false;
		dispatch('close');
	}

	let action: ModalAction = {
		text: 'Resize',
		isDisabled: true,
		func: async () => {
			dismissible = false;

			let res = await resizeMachine(machine.uuid, state.size);

			if (res.success) {
				return new Promise((resolve) => {
					showJobProgress(jobStatusEl, res.data.jobId, () => {
						dismissible = true;
						setTimeout(() => {
							dispatch('refresh');
							closeModal();
							resolve();
						}, 3000);
					}, () => {
						dismissible = true;
						setTimeout(() => {
							dispatch('refresh');
							resolve();
						}, 2000);
					})
			   });
			}
		}
	}

	$: action.isDisabled = machine && (state.size <= machine.size);  
</script>

<Modal title="Resize machine" type="warning" {action} {dismissible} {show} on:close={closeModal}>	
	<div class="mt-2 flex flex-col space-y-1">
		<div class="flex flex-col">
			<div class="font-bold text-sm tracking-wide opacity-50">Current size</div>
			<div class="text-xl -mt-1">{machine.size} GB</div>
		</div>
		<div class="flex flex-col pt-2">
			<div class="font-bold text-sm tracking-wide opacity-50">New size</div>
			<TextField
				on:change={(e) => (state.size = e.detail.value)}
				suffix="GB"
				type="number"
			/>
		</div>
		<div class="pt-4">
			<div bind:this={jobStatusEl} />
		</div>
	</div>
</Modal>