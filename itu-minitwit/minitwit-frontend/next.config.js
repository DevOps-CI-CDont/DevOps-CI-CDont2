/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  env: {
    NEXT_PUBLIC_PROXY_URL: "http://localhost:3001",
    NEXT_PUBLIC_API_URL: "http://localhost:8080"
  },
  output:'standalone'
}

module.exports = nextConfig
