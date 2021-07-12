<script lang="typescript">
	import { goto } from '$app/navigation';

	import { Button } from '../components/base';
	import { Steps, StepItem } from '../components/steps';
	import { SelectField, TextField } from '../components/fields';
	import { LatencyTest } from '../components/modals';

	import { completeSetup, verifyCredentials } from '../service/setup';

	const submit = async (): void => {
		stepper.nextStep();
	};

	const complete = async (): void => {
		let res = await completeSetup(
			state.form.awsAccessKeyId.value,
			state.form.awsSecretKey.value,
			state.form.region,
			state.form.password
		);

		if (res.success) {
			submit();
		} else {
			state.error.completeSetup = res.message;
		}
	};

	const verify = async (): void => {
		state.loading = true;

		let res = await verifyCredentials(
			state.form.awsAccessKeyId.value,
			state.form.awsSecretKey.value
		);

		state.loading = false;

		if (res.success) {
			state.regions = res.data.regions.map((r) => {
				return {
					text: r.name,
					value: r.id
				};
			});
			state.form.awsAccessKeyId.error = '';
			state.form.awsSecretKey.error = '';

			state.verifiedCreds = true;
		} else {
			state.form.awsAccessKeyId.error = 'Invalid credentials';
			state.form.awsSecretKey.error = 'Invalid credentials';
		}
	};

	let stepper;

	let state = {
		error: {
			completeSetup: ''
		},
		form: {
			awsAccessKeyId: {
				value: '',
				error: ''
			},
			awsSecretKey: {
				value: '',
				error: ''
			},
			password: '',
			confirmPassword: '',
			region: 'us-east-1'
		},
		modal: {
			show: false
		},
		regions: [],
		loading: false,
		verifiedCreds: false
	};

	$: disabled = {
		creds: state.form.awsAccessKeyId.value === '' || state.form.awsSecretKey.value === '',
		pwd: state.form.password !== state.form.confirmPassword || !state.form.password
	};
</script>

<svelte:head>
	<title>Setup</title>
</svelte:head>

<LatencyTest
	show={state.modal.show}
	on:close={() => (state.modal.show = false)}
	regions={state.regions}
/>

<div class="flex flex-col items-center w-full">
	<div class="min-w-full max-w-screen-lg h-full my-8">
		<div class="flex flex-col h-full justify-center items-center">
			<Steps bind:this={stepper}>
				<StepItem label="Begin">
					<div class="text-center h-64">
						<h1 class="text-8xl">Hi!</h1>
						<h2 class="text-3xl">I'm lakitu, your cloud gaming assistant</h2>
						<div class="mt-8">
							<Button text="Next" type="primary" size="lg" onClick={stepper.nextStep} />
						</div>
					</div>
				</StepItem>
				<StepItem label="Credentials">
					<div class="text-center h-64">
						<h1 class="text-3xl">Enter your credentials</h1>
						<div class="mt-8 md:w-1/2 mx-auto">
							<TextField
								on:change={(e) => (state.form.awsAccessKeyId.value = e.detail.value)}
								disabled={state.loading || state.verifiedCreds}
								error={state.form.awsAccessKeyId.error}
								type="password"
								placeholder="AWS Access Key ID"
								align="center"
								full
							/>
							<div class="my-4">
								<TextField
									on:change={(e) => (state.form.awsSecretKey.value = e.detail.value)}
									disabled={state.loading || state.verifiedCreds}
									error={state.form.awsSecretKey.error}
									type="password"
									placeholder="AWS Secret Key"
									align="center"
									full
								/>
							</div>
							<div class="flex justify-around space-x-2">
								<Button
									text="Verify"
									type="primary"
									disabled={disabled.creds || state.verifiedCreds}
									size="lg"
									onClick={verify}
									loading={state.loading}
									icon={state.verifiedCreds ? 'check-circle' : ''}
								/>
								<Button
									text="Next"
									type="primary"
									disabled={!state.verifiedCreds}
									size="lg"
									onClick={stepper.nextStep}
								/>
							</div>
						</div>
					</div>
				</StepItem>
				<StepItem label="Region">
					<div class="text-center h-64">
						<h1 class="text-3xl">Select a default AWS region</h1>
						<h2 class="text-xl">
							<span class="opacity-75">Not sure? Click</span>
							<a on:click={() => (state.modal.show = true)}>here</a>
							<span class="opacity-75">for a latency test.</span>
						</h2>
						<div class="mt-8 md:w-1/2 mx-auto">
							<div class="mb-4">
								<SelectField
									on:change={(e) => (state.form.region = e.detail.value)}
									placeholder="AWS Region"
									selectedValue={'us-east-1'}
									options={state.regions}
									full
								/>
							</div>
							<Button text="Next" type="primary" size="lg" onClick={stepper.nextStep} />
						</div>
					</div>
				</StepItem>
				<StepItem label="Security">
					<div class="text-center h-64">
						<h1 class="text-3xl">Create a password</h1>
						<div class="mt-8 md:w-1/2 mx-auto">
							<TextField
								on:change={(e) => (state.form.password = e.detail.value)}
								type="password"
								placeholder="Password"
								align="center"
								full
							/>
							<div class="my-4">
								<TextField
									on:change={(e) => (state.form.confirmPassword = e.detail.value)}
									type="password"
									placeholder="Confirm password"
									align="center"
									full
								/>
							</div>
							{#if state.error.completeSetup}{/if}
							<Button
								text="Next"
								type="primary"
								size="lg"
								disabled={disabled.pwd}
								onClick={complete}
							/>
						</div>
					</div>
				</StepItem>
				<StepItem label="Finish">
					<div class="text-center h-64">
						<h1 class="text-3xl">Setup complete!</h1>
						<h2 class="text-xl opacity-75">Login to create a cloud machine.</h2>
						<div class="mt-8">
							<Button text="Login" type="primary" size="lg" onClick={() => goto('/login')} />
						</div>
					</div>
				</StepItem>
			</Steps>
		</div>
	</div>
</div>
