/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  env: {
    NEXT_PUBLIC_PROXY_URL: "http://138.68.93.147:3001",
    NEXT_PUBLIC_API_URL: "http://138.68.93.147:8080",
  },
  output: "standalone",
};

module.exports = nextConfig;
