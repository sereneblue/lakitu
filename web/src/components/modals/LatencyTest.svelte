<script lang="typescript">
	import { Button, Modal } from '../base';
	import { pingAWS } from '../../service/setup';

	export let regions: string[] = [];
	export let show: boolean = false;

	type PingResult = {
		region: string;
		latency: number;
	};

	const testLatency = async (): Promise<void> => {
		loading = true;
		pingResults = [];

		let res = await pingAWS(regions.map((r) => r.value));

		if (res.success) {
			for (let i = 0; i < regions.length; i++) {
				pingResults.push({
					name: regions[i].text,
					latency: res.data.latency[regions[i].value]
				});
			}
			pingResults.sort((a, b) => a.latency - b.latency);
			pingResults = pingResults;
		}

		loading = false;
	};

	let pingResults: PingResult[] = [];
	let loading: boolean = false;
</script>

<Modal title="AWS Latency Test" type="info" {show} on:close={() => (show = false)}>
	{#if pingResults.length}
		<table class="w-full text-left divide-y divide-gray-200 block max-h-60 overflow-y-auto">
			<thead class="bg-gray-50 block">
				<tr>
					<th scope="col" class="px-6 py-2 text-xs font-medium text-gray-500 uppercase  w-full">
						Region
					</th>
					<th scope="col" class="px-6 py-2 text-xs font-medium text-gray-500 uppercase">
						Latency
					</th>
				</tr>
			</thead>
			<tbody class="bg-white divide-y divide-gray-200 text-xs text-gray-900 block">
				{#each pingResults as p}
					<tr>
						<td class="px-6 py-1 w-full">
							{p.name}
						</td>
						<td class="px-6 whitespace-nowrap py-1">
							{p.latency} ms
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	{/if}
	<div class="flex justify-center mt-4">
		<Button type="primary" text="Test" onClick={testLatency} {loading} />
	</div>
</Modal>
