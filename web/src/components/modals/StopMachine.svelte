<script lang="ts">
	import { createEventDispatcher } from 'svelte'; 
	import { Modal } from '../base';
	import { TextField } from '../fields';

	import type { MachineData, ModalAction } from '../../service/types';
	import { stopMachine } from '../../service/machines';
	import { showJobProgress } from '../../service/util';

	export let show: boolean = false;
	export let machine: MachineData;

	let dismissible = true;
	let jobStatusEl: HTMLElement;

	const dispatch = createEventDispatcher();

	const closeModal = () => {
		show = false;
		dispatch('close');
	}

	let action: ModalAction = {
		text: 'Yes, stop!',
		isDisabled: false,
		func: async () => {
			dismissible = false;

			let res = await stopMachine(machine.uuid);

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
</script>

<Modal title="Stop machine" type="danger" {action} {dismissible} {show} on:close={closeModal}>	
	<div class="mt-2 flex flex-col space-y-1">
		<div class="mt-2">
			<h2 class="font-semibold text-lg">
				Are you sure you want to stop this machine?
			</h2>
		</div>
		<div class="pt-4">
			<div bind:this={jobStatusEl} />
		</div>
	</div>
</Modal>