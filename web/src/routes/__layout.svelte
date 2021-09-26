<div>
	<div class="flex bg-nord0 text-nord4 h-full min-h-screen px-6 py-8">
		<slot />
	</div>
</div>

<style lang="postcss" global>
	@tailwind base;
	@tailwind components;
	@tailwind utilities;

	a {
		@apply cursor-pointer
			   underline
			   text-accent-200 hover:text-accent-300 font-semibold;
	}
</style>
