<script>
	import { addMessages, getLocaleFromNavigator, init } from 'svelte-i18n/dist/runtime.cjs.js';

	import en from '../locales/en.json';

	addMessages('en', en);
	addMessages('en-US', en);

	init({
		fallbackLocale: 'en',
		initialLocale: getLocaleFromNavigator()
	});
</script>

<div class="dark">
	<div class="flex bg-nord6 text-nord2 dark:bg-nord0 dark:text-nord4 h-full min-h-screen px-6 py-8">
		<slot />
	</div>
</div>

<style global>
	@tailwind base;
	@tailwind components;
	@tailwind utilities;
</style>
