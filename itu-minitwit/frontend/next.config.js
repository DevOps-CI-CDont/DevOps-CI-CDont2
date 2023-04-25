/** @type {import('next').NextConfig} */
const nextConfig = {
	reactStrictMode: true,
	output: "standalone",
	typescript: {
		ignoreBuildErrors: false,
	},
	eslint: {
		ignoreDuringBuilds: true,
	},
	async redirects() {
		return [
		  {
			source: '/grafana',
			destination: ':4000',
			permanent: true,
		  },
		]
	},
};

module.exports = nextConfig;
