import { addMessages, getLocaleFromNavigator, init } from 'svelte-i18n/dist/runtime.cjs.js';

import en from './locales/en.json';

addMessages('en', en);
addMessages('en-US', en);

init({
  fallbackLocale: 'en',
  initialLocale: getLocaleFromNavigator()
});