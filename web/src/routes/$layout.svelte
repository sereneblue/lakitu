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

<slot />

<style global>
	@tailwind base;
	@tailwind components;
	@tailwind utilities;
</style>
