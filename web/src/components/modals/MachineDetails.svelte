<script lang="ts">
	import { createEventDispatcher } from 'svelte'; 
	import { Modal } from '../base';

	import type { MachineData } from '../../service/types';

	export let show: boolean = false;
	export let machine: MachineData;

	const dispatch = createEventDispatcher();

	const closeModal = () => {
		show = false;
		dispatch('close');
	}
</script>

<Modal title="Machine details" type="info" {show} on:close={closeModal}>
	<div class="my-2 flex flex-col space-y-1">
		<div class="flex flex-col">
			<div class="font-bold text-xs tracking-wide opacity-50">Name</div>
			<div class="text-lg -mt-1">{machine.name}</div>
		</div>
		<div class="flex flex-col">
			<div class="font-bold text-xs tracking-wide opacity-50">AMI ID</div>
			<div class="text-lg -mt-1">{machine.amiId}</div>
		</div>
		<div class="flex flex-col">
			<div class="font-bold text-xs tracking-wide opacity-50">Snapshot ID</div>
			<div class="text-lg -mt-1">{machine.snapshotId}</div>
		</div>
		<div class="flex flex-col">
			<div class="font-bold text-xs tracking-wide opacity-50">Password</div>
			<div class="text-lg -mt-1">{machine.adminPwd}</div>
		</div>
	</div>
</Modal>