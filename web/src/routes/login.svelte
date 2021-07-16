<script lang="typescript">
	import { goto } from '$app/navigation';
	import { _ } from 'svelte-i18n/dist/runtime.cjs.js';

	import { Button, Icon, Modal } from '../components/base';
	import { TextField } from '../components/fields';
	import { login } from '../service/session';

	let state = {
		password: {
			value: '',
			error: ''
		}
	};

	const handleLogin = async (): Promise<void> => {
		let res = await login(state.password.value);

		if (res.success) {
			goto('/machines');
		} else {
			state.password.error = res.message;
		}
	};
</script>

<svelte:head>
	<title>Login</title>
</svelte:head>

<div class="flex flex-col items-center w-full">
	<div class="min-w-full max-w-screen-lg h-full my-8">
		<div class="flex flex-col h-full justify-center items-center">
			<div class="text-6xl">lakitu</div>
			<div class="w-4/5 md:w-2/5 my-8">
				<TextField
					on:change={(e) => (state.password.value = e.detail.value)}
					type="password"
					placeholder={$_('fields.password')}
					align="center"
					error={state.password.error}
					submit={handleLogin}
					full
				/>
			</div>
			<div>
				<Button type="primary" text={$_('buttons.login')} size="lg" onClick={handleLogin} />
			</div>
		</div>
	</div>
</div>
