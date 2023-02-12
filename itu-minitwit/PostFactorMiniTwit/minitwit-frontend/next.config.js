/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
}

module.exports = {
  nextConfig,
  env: {
    NEXT_PUBLIC_CORS_ORIGIN: "http://localhost:6001",
    NEXT_PUBLIC_API_URL: "http://127.0.0.1:8080"
  }
}

// module.exports = nextConfig
