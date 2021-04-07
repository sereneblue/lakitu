module.exports = {
	purge: {
		mode: 'all',
		content: ['./src/**/*.html', './src/**/*.svelte'],
		options: {
			whitelistPatterns: [/svelte-/],
			defaultExtractor: (content) =>
				[...content.matchAll(/(?:class:)*([\w\d-/:%.]+)/gm)].map(
					([_match, group, ..._rest]) => group
				),
			keyframes: true
		}
	},
	mode: 'jit',
	darkMode: 'class',
	theme: {
		extend: {}
	},
	variants: {
		extend: {}
	},
	plugins: []
};
