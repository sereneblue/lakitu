<script lang="ts">
	import { createEventDispatcher } from "svelte";
	import { fade } from 'svelte/transition';

	import { Icon } from '../base';
	import type { MachineEvent } from '../../service/types';

	export let status: MachineStatus;
	export let open: boolean = false;

	const dispatch = createEventDispatcher();
	
	const emit = (event: MachineEventType): void => {
		dispatch("emit", event);
		open = false;
	}

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
		return clickOutside(node, () => {
			open = false;
		});
	};
</script>

{#if open}
	<div class="flex absolute top-0 right-0 mt-2 mr-2 cursor-pointer">
		<div class="relative inline-block text-left mr-4">
			<div
			    transition:fade={{ duration: 150 }}
			    use:handleClickOutside
			    class="origin-top-right absolute right-0 w-40 z-10 rounded shadow-lg">
			    <div class="rounded border overflow-hidden border-indigo-500 bg-white text-nord2">
			    	{#if status === "offline" }
						<div
					        class="flex hover:bg-nord2/10"
					        role="menu"
					        aria-orientation="horizontal"
					        aria-labelledby="options-menu">
					        <button
					          on:click={(e) => emit('start')}
					          class="flex items-center w-full p-1.5 text-sm"
					          role="menuitem">
					          <span class="w-3 h-3 mx-1.5">
					          	<Icon icon="play" />
					          </span>
					          Start
					        </button>
						</div>
						<div
					        class="flex hover:bg-nord2/10"
					        role="menu"
					        aria-orientation="horizontal"
					        aria-labelledby="options-menu">
					        <button
					          on:click={(e) => emit('transfer')}
					          class="flex items-center w-full p-1.5 text-sm"
					          role="menuitem">
					          <span class="w-3 h-3 mx-1.5">
					          	<Icon icon="refresh-cw" />
					          </span>
					          Change region
					        </button>
						</div>
						<div
					        class="flex hover:bg-nord2/10"
					        role="menu"
					        aria-orientation="horizontal"
					        aria-labelledby="options-menu">
					        <button
					          on:click={(e) => emit('resize')}
					          class="flex items-center w-full p-1.5 text-sm"
					          role="menuitem">
					          <span class="w-3 h-3 mx-1.5">
					          	<Icon icon="plus" />
					          </span>
					          Add storage
					        </button>
						</div>
					{:else if status === "online"}
						<div
					        class="flex hover:bg-nord2/10"
					        role="menu"
					        aria-orientation="horizontal"
					        aria-labelledby="options-menu">
					        <button
					          on:click={(e) => emit('stop')}
					          class="flex items-center w-full p-1.5 text-sm"
					          role="menuitem">
					          <span class="w-3 h-3 mx-1.5">
					          	<Icon icon="square" />
					          </span>
					          Stop
					        </button>
						</div>
					{/if}
					<div
				        class="flex hover:bg-nord2/10"
				        role="menu"
				        aria-orientation="horizontal"
				        aria-labelledby="options-menu">
				        <button
				          on:click={(e) => emit('delete')}
				          class="flex items-center w-full p-1.5 text-sm"
				          role="menuitem">
				          <span class="w-3 h-3 mx-1.5">
				          	<Icon icon="trash" />
				          </span>
				          Delete
				        </button>
					</div>
			    </div>
			</div>
		</div>
	</div>
{/if}