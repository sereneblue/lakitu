<script lang="ts">
	import { onMount } from "svelte";
	import { goto } from '$app/navigation';
	import { getJobStatus } from "../service/job";
	import { notify } from "../service/util";

	let state = {
		job: {},
		progress: 0,
	};

	let errorEl: HTMLElement;

	$: state.progress = ( state.job.completed / state.job.total ) * 100;

	const checkJob = async (): Promise<void> => {
		let res = await getJobStatus();
		
		if (res.success) {
			if (res.data == undefined || (res.data && res.data.isComplete)) {
				window.location.href = '/machines';
			} else {
				state.job = res.data;		
			}
		} else {
			notify(errorEl, "danger", res.message, 0);
		}

		setTimeout(checkJob, 2000);
	};
 
	onMount(() => {
		checkJob();
	})
</script>

<svelte:head>
	<title>Loading...</title>
</svelte:head>

<div class="w-full h-auto">
	<div class="flex flex-col h-full justify-center items-center">
		{#if state.job.name }
			<div class="block overflow-hidden w-full max-w-lg mb-4 h-1.5 rounded bg-white/25 shadow">
				<div style="width: { state.progress }%" class="h-1.5 bg-accent-100 transition-all duration-500"></div>
			</div>
			<div class="text-center">
				<h1 class="text-lg font-bold">
					{state.job.name}
				</h1>
				<h2 class="text-base text-opacity-75">
					{state.job.currentTask}
				</h2>
				<div class="mt-4">
					{state.job.completed } / {state.job.total}
				</div>
			</div>
		{/if}
		<div bind:this={errorEl} class="mt-2"></div>
	</div>
</div>
