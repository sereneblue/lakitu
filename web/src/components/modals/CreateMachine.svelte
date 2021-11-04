<script lang="ts">
	import { createEventDispatcher, onMount, onDestroy } from 'svelte';

	import { Button, Icon, Modal } from '../base';
	import { TextField, SelectField } from '../fields';

	import { getInstances, getPricing, createMachine } from '../../service/machines';
	import { showJobProgress } from '../../service/util';

	export let defaultRegion: string = 'us-east-1';
	export let regions: string[] = [];
	export let show: boolean = false;

	let dismissible = true;
	let jobStatusEl: HTMLElement;

	const dispatch = createEventDispatcher();

	const closeModal = () => {
		show = false;
		dispatch('close');
	}

	const create = async (): Promise<void> => {
		dismissible = false;

		let res = await createMachine(state.form);

		if (res.success) {
			showJobProgress(jobStatusEl, res.data.jobId, () => {
				dismissible = true;
				setTimeout(() => {
					dispatch('refresh');
					closeModal();
					resolve();
				}, 3000);
			})
		}
	}

	const updatePrices = async (region: string) => {
		// use cache to update price
		if (state.regions[region] && state.regions[region].pricing.volume) {
			state.form.region = region;
			state.form.instanceType = '';
			return;
		}

		state.regions[region] = {
			pricing: {},
			instances: {}
		};

		let res = await Promise.all([getPricing(region), getInstances(region)]);

		if (res[0].success) {
			state.regions[region].pricing = {};
			state.regions[region].pricing.spotInstance = 0;
			state.regions[region].pricing.volume = res[0].data.volume;
			state.regions[region].pricing.snapshot = res[0].data.snapshot;
			state.regions[region].pricing.bandwidth = res[0].data.bandwidth;
		}

		if (res[1].success) {
			state.regions[region].availableInstances = res[1].data.instances.map((i) => {
				return {
					text: i.instance,
					value: i.instance
				};
			});

			for (let i = 0; i < res[1].data.instances.length; i++) {
				state.regions[region].instances[res[1].data.instances[i].instance] =
					res[1].data.instances[i].price;
			}
		}

		state.form.region = region;
		state.form.instanceType = '';
	};

	let state = {
		regions: {
			[defaultRegion]: {
				availableInstances: [],
				pricing: {
					spotInstance: 0,
					volume: 0,
					bandwidth: 0,
					snapshot: 0
				}
			}
		},
		numHours: 10,
		form: {
			name: '',
			instanceType: '',
			region: defaultRegion,
			size: 50,
			streamOption: 'parsec'
		},
		streamOptions: [{ text: 'Parsec', value: 'parsec' }]
	};

	$: cost = {
		storage:
			state.form.size *
				((730 - state.numHours) / 730) *
				state.regions[state.form.region].pricing.snapshot +
			state.numHours * state.regions[state.form.region].pricing.volume,
		network: state.numHours * state.regions[state.form.region].pricing.bandwidth * 10,
		instance: state.numHours * state.regions[state.form.region].pricing.spotInstance
	};

	$: totalCost = (cost.storage + cost.network + cost.instance).toFixed(2);

	$: validForm =
		state.form.name &&
		state.form.streamOption &&
		state.form.region &&
		state.form.instanceType &&
		state.form.size;

	$: if (!show) {
		state.form = {
			name: '',
			instanceType: '',
			region: defaultRegion,
			size: 50,
			streamOption: 'parsec'
		};

		state.numHours = '';
	}

	// watch gpu instance
	$: {
		if (state.form.instanceType != '') {
			state.regions[state.form.region].pricing.spotInstance = state.regions[state.form.region].instances[state.form.instanceType];
		} else {
			state.regions[state.form.region].pricing.spotInstance = 0;
		}
	}

	onMount(async () => {
		await updatePrices(state.form.region);
	});
</script>

<Modal title="Create a machine" {dismissible} {show} on:close={closeModal}>
	<div class="mt-2">
		<TextField on:change={(e) => (state.form.name = e.detail.value.trim())} label="Name" full />
	</div>
	<div class="mt-2 flex space-x-4">
		<div class="w-1/2">
			<SelectField
				on:change={(e) => updatePrices(e.detail.value)}
				label="Region"
				selectedValue={state.form.region}
				options={regions}
				full
			/>
		</div>
		<div class="w-1/2">
			<SelectField
				on:change={(e) => (state.form.streamOption = e.detail.value)}
				label="Stream software"
				selectedValue={state.form.streamOption}
				options={state.streamOptions}
				full
			/>
		</div>
	</div>
	<div class="mt-2 flex space-x-4">
		<div class="w-1/2">
			<SelectField
				on:change={(e) => (state.form.instanceType = e.detail.value)}
				label="Instance Type"
				selectedValue={state.form.instanceType}
				placeholder="Select instance"
				options={state.regions[state.form.region].availableInstances}
				full
			/>
		</div>
		<div class="w-1/2">
			<TextField
				on:change={(e) => (state.form.size = e.detail.value)}
				value={state.form.size}
				label="Storage"
				suffix="GB"
				type="number"
				full
			/>
		</div>
	</div>
	<div class="mt-4 border-t pt-2">
		<div>Estimated Cost</div>
		<div class="flex justify-between items-center">
			<TextField
				on:change={(e) => (state.numHours = e.detail.value)}
				placeholder="Hours per month"
				value={state.numHours}
				type="number"
			/>
			<div class="font-bold text-lg flex justify-between items-center">
				<span class="mr-1">≈ ${totalCost}/month</span>
				<span class="peer w-4 h-4 cursor-pointer"><Icon icon="help-circle" /></span>
				<div
					class="font-normal hidden peer-hover:block text-sm absolute right-6 p-2 -mt-32 bg-white shadow-lg border rounded"
				>
					<table>
						<tr>
							<td>Storage</td>
							<td>${cost.storage.toFixed(2)}</td>
						</tr>
						<tr>
							<td>Network (10MBPS)</td>
							<td>${cost.network.toFixed(2)}</td>
						</tr>
						<tr>
							<td>Instance</td>
							<td>${cost.instance.toFixed(2)}</td>
						</tr>
					</table>
				</div>
			</div>
		</div>
	</div>
	<div class="pt-4">
		<div bind:this={jobStatusEl} />
	</div>
	<div class="flex justify-center mt-4">
		<Button type="primary" text="Create" disabled={!validForm} onClick={create} />
	</div>
</Modal>
