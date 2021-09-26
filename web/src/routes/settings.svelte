<script context="module">
	export async function load({ page, fetch, session, context }) {
		let res = await getUserData(page.host);

		if (res.success) {
			return {
				props: {
					userData: res.data
				}
			};
		}

		return true;
	}
</script>

<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';

	import { Button, Icon, Navigation } from '../components/base';
	import { TextField, SelectField } from '../components/fields';
	import { LatencyTest } from '../components/modals';

	import { getUserData, updatePassword, updatePreferences } from '../service/session';
	import { notify } from '../service/util';

	enum ACTIVE_TAB {
		AWS,
		USER,
		ABOUT
	}

	const handleTabClick = (newTab: ACTIVE_TAB): void => {
		state.activeTab = newTab;

		switch (newTab) {
			case ACTIVE_TAB.USER:
				location.hash = '#user';
				break;
			case ACTIVE_TAB.ABOUT:
				location.hash = '#about';
				break;
			default:
				location.hash = '#aws';
				break;
		}
	};

	export let userData: object = {
		regions: [],
		accessKey: '',
		secretKey: '',
		defaultRegion: ''
	};

	const handleUpdatePassword = async (): Promise<void> => {
		let res = await updatePassword(
			state.form.password,
			state.form.newPassword,
			state.form.confirmNewPassword
		);

		notify(state.notifications.updatePassword, res.success ? 'success' : 'danger', res.message);
	};

	const handleUpdatePreferences = async (): Promise<void> => {
		let res = await updatePreferences(
			state.form.awsAccessKeyId,
			state.form.awsSecretKey,
			state.form.defaultRegion
		);

		notify(state.notifications.updatePreferences, res.success ? 'success' : 'danger', res.message);
	};

	onMount(() => {
		switch (location.hash.toLowerCase()) {
			case '#user':
				state.activeTab = ACTIVE_TAB.USER;
				break;
			case '#about':
				state.activeTab = ACTIVE_TAB.ABOUT;
				break;
			default:
				state.activeTab = ACTIVE_TAB.AWS;
				break;
		}

		state.regions = userData.regions.map((r) => {
			return {
				text: r.name,
				value: r.id
			};
		});
		state.form.awsAccessKeyId = userData.accessKey;
		state.form.awsSecretKey = userData.secretKey;
		state.form.defaultRegion = userData.defaultRegion;
	});

	let state = {
		form: {
			awsAccessKeyId: '',
			awsSecretKey: '',
			defaultRegion: 'us-east-1',
			newPassword: '',
			confirmNewPassword: '',
			password: ''
		},
		activeTab: null,
		notifications: {
			updatePassword: null,
			updatePreferences: null
		},
		regions: [],
		showLatencyModal: false,
		settings: [
			{ text: 'AWS', value: ACTIVE_TAB.AWS },
			{ text: 'User', value: ACTIVE_TAB.USER },
			{ text: 'About', value: ACTIVE_TAB.ABOUT }
		]
	};

	$: pwdDisabled =
		state.form.newPassword !== state.form.confirmNewPassword ||
		!state.form.newPassword ||
		!state.form.password;

	$: prefDisabled = !state.form.awsAccessKeyId || !state.form.awsSecretKey;
</script>

<svelte:head>
	<title>Settings</title>
</svelte:head>

<Navigation />

<LatencyTest
	show={state.showLatencyModal}
	on:close={() => (state.showLatencyModal = false)}
	regions={state.regions}
/>

<div class="flex flex-col items-center w-full" id="settings">
	<div class="min-w-full max-w-screen-lg h-full my-4">
		<div class="max-w-screen-xl mx-auto">
			<div class="mt-8">
				<div class="flex justify-between items-center mb-4">
					<h2 class="text-2xl mb-4">Settings</h2>
				</div>
				<div class="flex flex-col sm:flex-row">
					<ul class="w-1/3 hidden sm:block space-y-2 mr-4">
						<li
							class="p-2 cursor-pointer rounded flex space-x-2 items-center"
							class:active={state.activeTab === ACTIVE_TAB.AWS}
							on:click={(e) => handleTabClick(ACTIVE_TAB.AWS)}
						>
							<div class="w-5 h-5">
								<Icon icon="cloud" />
							</div>
							<span> AWS </span>
						</li>
						<li
							class="p-2 cursor-pointer rounded flex space-x-2 items-center"
							class:active={state.activeTab === ACTIVE_TAB.USER}
							on:click={(e) => handleTabClick(ACTIVE_TAB.USER)}
						>
							<div class="w-5 h-5">
								<Icon icon="user" />
							</div>
							<span> User </span>
						</li>
						<li
							class="p-2 cursor-pointer rounded flex space-x-2 items-center"
							class:active={state.activeTab === ACTIVE_TAB.ABOUT}
							on:click={(e) => handleTabClick(ACTIVE_TAB.ABOUT)}
						>
							<div class="w-5 h-5">
								<Icon icon="alert-circle" />
							</div>
							<span> About </span>
						</li>
					</ul>
					<div class="w-full block sm:hidden mb-4">
						<SelectField
							on:change={(e) => (state.activeTab = Number(e.detail.value))}
							selectedValue={state.activeTab}
							options={state.settings}
							full
						/>
					</div>
					<div class="w-full">
						<div class="flex flex-col max-w-md space-y-4">
							{#if state.activeTab === ACTIVE_TAB.AWS}
								<div>
									<TextField
										on:change={(e) => (state.form.awsAccessKeyId = e.detail.value.trim())}
										value={state.form.awsAccessKeyId}
										label="Access Key"
										type="password"
										full
									/>
								</div>
								<div>
									<TextField
										on:change={(e) => (state.form.awsSecretKey = e.detail.value.trim())}
										value={state.form.awsSecretKey}
										type="password"
										label="Secret Key"
										full
									/>
								</div>
								<div class="flex space-x-2 items-end">
									<div class="flex-1">
										<SelectField
											on:change={(e) => (state.form.defaultRegion = e.detail.value)}
											label="Default Region"
											selectedValue={state.form.defaultRegion}
											options={state.regions}
											full
										/>
									</div>
									<Button
										type="primary"
										icon="radio"
										disabled={prefDisabled}
										onClick={() => (state.showLatencyModal = true)}
									/>
								</div>
								<Button
									type="primary"
									text="Update preferences"
									disabled={prefDisabled}
									onClick={handleUpdatePreferences}
								/>
								<div bind:this={state.notifications.updatePreferences} />
							{:else if state.activeTab === ACTIVE_TAB.USER}
								<div>
									<TextField
										type="password"
										on:change={(e) => (state.form.newPassword = e.detail.value.trim())}
										label="New Password"
										full
									/>
								</div>
								<div>
									<TextField
										type="password"
										on:change={(e) => (state.form.confirmNewPassword = e.detail.value.trim())}
										label="Confirm New Password"
										full
									/>
								</div>
								<div>
									<TextField
										type="password"
										on:change={(e) => (state.form.password = e.detail.value.trim())}
										label="Currrent Password"
										full
									/>
								</div>
								<Button
									type="primary"
									text="Change password"
									disabled={pwdDisabled}
									onClick={handleUpdatePassword}
								/>
								<div bind:this={state.notifications.updatePassword} />
							{:else if state.activeTab === ACTIVE_TAB.ABOUT}
								<div class="w-full text-lg">
									<div>v0.1 by sereneblue</div>
									<div class="mt-4">
										<a
											href="https://github.com/sereneblue/lakitu"
											target="_blank"
											rel="noopener"
											class="flex items-center space-x-2 no-underline"
										>
											<span class="w-5 h-5"><Icon icon="github" /></span>
											<span>Source Code</span>
										</a>
										<a
											href="https://github.com/sereneblue/lakitu/issues"
											target="_blank"
											rel="noopener"
											class="flex items-center space-x-2 no-underline"
										>
											<span class="w-5 h-5"><Icon icon="alert-triangle" /></span>
											<span>Report a bug</span>
										</a>
									</div>
								</div>
							{/if}
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</div>

<style lang="postcss" global>
	#settings .active {
		@apply bg-nord4 dark:bg-nord1 font-bold text-accent-100;
	}
	#settings li:not(.active) {
		@apply opacity-50 hover:bg-nord4/75 dark:hover:bg-nord1/75;
	}
</style>
