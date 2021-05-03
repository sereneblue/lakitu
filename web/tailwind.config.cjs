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
		extend: {
			colors: {
				nord0: '#2e3440',
				nord1: '#3b4252',
				nord2: '#434c5e',
				nord3: '#4c566a',
				nord4: '#d8dee9',
				nord5: '#e5e9f0',
				nord6: '#eceff4',
				accent: {
					100: '#81a9a8',
					200: '#729696',
					300: '#648483'
				},
				danger: {
					100: '#bf616a',
					200: '#994e55',
					300: '#733a40'
				},
				info: '#5e81ac',
				success: '#a3be8c',
				warning: '#ebcb8b'
			}
		}
	},
	variants: {
		extend: {}
	},
	plugins: []
};
