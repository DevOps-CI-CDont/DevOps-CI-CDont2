/** @type {import('next').NextConfig} */
const nextConfig = {
	reactStrictMode: true,
	env: {
		NEXT_PUBLIC_URL: "http://138.68.93.147:3000/public",
		NEXT_PUBLIC_API_URL: "http://138.68.93.147:8080",
	},
	output: "standalone",
	typescript: {
		ignoreBuildErrors: true,
	},
	eslint: {
		ignoreDuringBuilds: true,
	},
};

module.exports = nextConfig;
