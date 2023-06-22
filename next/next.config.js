const { remarkCodeHike } = require("@code-hike/mdx");
const theme = require("shiki/themes/monokai.json");

/** @type {import('next').NextConfig} */
const nextConfig = {
  // Configure pageExtensions to include md and mdx
  pageExtensions: ['ts', 'tsx', 'js', 'jsx', 'md', 'mdx'],
  // Optionally, add any other Next.js config below
  reactStrictMode: true,

	typescript: {
		ignoreBuildErrors: true,
	},

	images: {
		remotePatterns: [{
			protocol: "https",
			hostname: "lh3.googleusercontent.com",
			port: "",
			pathname: "/**",
		}]
	}
}

const withMDX = require('@next/mdx')({
  extension: /\.mdx?$/,
  options: {
    remarkPlugins: [[remarkCodeHike, { 
			showCopyButton: true,
			theme,
		}]],
    rehypePlugins: [],
    providerImportSource: "@mdx-js/react",
  },
})

module.exports = withMDX(nextConfig)
