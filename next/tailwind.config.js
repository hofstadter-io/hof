/** @type {import('tailwindcss').Config} */
module.exports = {
	darkMode: 'class',
  content: [
    './src/pages/**/*.{js,ts,jsx,tsx,mdx}',
    './src/components/**/*.{js,ts,jsx,tsx,mdx}',
    './src/app/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  plugins: [
		require('@tailwindcss/typography'),
		require("daisyui"),
	],

	daisyui: {
		logs: false,
		// styled: false,
	},

  theme: {
		extend: {
			textDecorationThickness: {
				3: '3px',
			}
		}
	}
}
