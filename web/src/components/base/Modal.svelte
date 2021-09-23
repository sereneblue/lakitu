<script lang="typescript">
	import { createEventDispatcher } from 'svelte';
	import { fade, fly } from 'svelte/transition';

	import Icon from './Icon.svelte';
	import Button from './Button.svelte';
	import type { ModalType, ModalAction } from '../../service/types';

	export let type: ModalType = 'default';
	export let action: ModalAction;
	export let title: string = '';
	export let text: string = '';
	export let dismissible: boolean = true;
	export let show: boolean = false;

	const closeModal = (): void => {
		if (dismissible) {
			show = false;
		}
	};

	const clickOutside = (node: Node, closeAction: Function): object => {
		const handleClick = (event: Event) => {
			let path = event.composedPath();

			if (!path.includes(node)) {
				closeAction();
			}
		};

		setTimeout(() => {
			document.addEventListener('click', handleClick);
		}, 10);

		return {
			destroy() {
				document.removeEventListener('click', handleClick);
			}
		};
	};

	const handleClickOutside = (node: Node): object => {
		return clickOutside(node, closeModal);
	};

	const init = (): object => {
		const handleClose = (e: KeyboardEvent) => {
			if (e.key === 'Escape' && dismissible) {
				closeModal();
			}
		};

		window.addEventListener('keydown', handleClose);

		return {
			destroy() {
				window.removeEventListener('keydown', handleClose);
			}
		};
	};

	const dispatch = createEventDispatcher();

	$: if (!show) dispatch('close');
</script>

{#if show}
	<div
		class="fixed z-20 inset-0 overflow-y-auto"
		aria-labelledby="modal-title"
		role="dialog"
		aria-modal="true"
		use:init
		transition:fade={{ duration: 250 }}
	>
		<div
			class="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0"
		>
			<div class="fixed inset-0 bg-nord0 bg-opacity-80 transition-opacity" aria-hidden="true" />
			<span class="hidden sm:inline-block sm:align-middle sm:h-screen" aria-hidden="true"
				>&#8203;</span
			>
			{#if show}
				<div
					use:handleClickOutside
					in:fly={{ y: 100, duration: 250 }}
					out:fly={{ y: 100, duration: 250 }}
					class="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform sm:my-8 sm:align-middle sm:max-w-lg w-full"
				>
					<div class="relative bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
						{#if dismissible}
							<div
								class="absolute top-0 right-0 h-6 w-6 m-4 cursor-pointer text-gray-900 opacity-25 hover:opacity-50"
								on:click={closeModal}
							>
								<Icon icon="x" />
							</div>
						{/if}
						<div class="sm:flex sm:items-start">
							{#if type != 'default'}
								<div class="w-full sm:w-12 flex sm:block justify-center">
									<div
										class="h-12 w-12 mr-4 flex-shrink-0 sm:mx-0 sm:h-10 sm:w-10"
										class:text-danger-200={type === 'danger'}
										class:text-info={type === 'info'}
										class:text-warning={type === 'warning'}
										class:text-success={type === 'success'}
									>
										{#if type === 'danger'}
											<Icon icon="alert-octagon" />
										{:else if type === 'info'}
											<Icon icon="info" />
										{:else if type === 'warning'}
											<Icon icon="alert-triangle" />
										{:else if type === 'success'}
											<Icon icon="check-circle" />
										{/if}
									</div>
								</div>
							{/if}
							<div
								class="mt-3 w-full sm:mt-0 text-left text-nord2"
								class:sm:ml-2={type != 'default'}
							>
								<h3 class="text-lg leading-6 font-medium text-center sm:text-left" id="modal-title">
									{title}
								</h3>
								<div class="mt-4 max-h-84 modal-content">
									{#if text}
										<p class="text-sm">
											{text}
										</p>
									{:else}
										<slot />
									{/if}
								</div>
							</div>
						</div>
					</div>
					{#if action}
						<div class="bg-gray-100 px-4 py-2 sm:px-6 sm:flex sm:flex-row-reverse">
							<span class="mx-2">
								<Button
									text={action.text}
									type={type === 'danger' ? 'danger' : 'primary'}
									onClick={action.func}
								/>
							</span>
							{#if !action.hideCancel}
								<Button text="Cancel" once outline flat onClick={closeModal} />
							{/if}
						</div>
					{/if}
				</div>
			{/if}
		</div>
	</div>
{/if}
