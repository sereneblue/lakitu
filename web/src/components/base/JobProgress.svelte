<script lang="ts">
	import { onMount } from "svelte";
	import { goto } from '$app/navigation';
	import { getJobStatus } from "../../service/job";
	import { notify } from "../../service/util";
	import { Status } from '../../service/types';

	export let jobId: number = 0;
	export let showAllDetails: boolean = false;
	export let done: Function = () => {};

	let state = {
		job: {
			name: '',
			current: '',
			status: 0,
			completed: 0,
			total: 0
		},
		id: 0,
		progress: 0,
	};

	let notifyEl: HTMLElement;
	let timeout;

	const checkJob = async (): Promise<void> => {
		let res = await getJobStatus(state.id);

		if (res.success) {
			state.job = res.data;

			if (res.data.isComplete) {
				notify(notifyEl, "success", "Task was completed successfully", 0);

				done();

				clearTimeout(timeout);
			} else {
				if (res.data.status == Status.ERROR) {
					notify(notifyEl, "danger", "There was an error completing task.", 0);
				} else {
					timeout = setTimeout(checkJob, 2000);
				}
			}
		} else {
			notify(notifyEl, "danger", res.message, 0);
			clearTimeout(timeout);
		}
	};

	$: state.progress = ( state.job.completed / state.job.total ) * 100;
 
	onMount(() => {
		state.id = jobId;
		checkJob();
	})	
</script>

<div class="w-full h-auto">
	<div class="flex flex-col h-full justify-center items-center">
		{#if state.id > 0}
			<div class="block overflow-hidden w-full max-w-lg mb-2 h-2 rounded bg-white/25 border border-1 border-nord0/50 shadow">
				<div style="width: {state.progress}%" class="h-2 bg-accent-100 transition-all duration-500 animate-pulse" />
			</div>
			<div class="text-center">
				{#if showAllDetails}
					<h1 class="text-lg font-bold">
						{state.job.name}
					</h1>
				{/if}
				<h2 class="text-base text-opacity-75">
					{state.job.currentTask} [ {state.job.completed } / {state.job.total} ]
				</h2>
			</div>
		{/if}
		<div bind:this={notifyEl} class="mt-2"></div>
	</div>
</div>