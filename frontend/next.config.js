/** @type {import('next').NextConfig} */
const nextConfig = {
  async rewrites() {
    return [
      {
        source: '/api/:path*',
        destination: 'http://backend:8080/:path*',
      },
    ]
  },
  async redirects() {
    return [
      {
        source: '/booking',
        destination: '/booking/selection',
        permanent: true,
      },
    ]
  },
}

module.exports = nextConfig
