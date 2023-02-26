/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  env: {
    NEXT_PUBLIC_PROXY_URL: "http://localhost:3001",
    NEXT_PUBLIC_API_URL: "https://seashell-app-hlfb2.ondigitalocean.app"
  },
  output:'standalone'
}

module.exports = nextConfig
