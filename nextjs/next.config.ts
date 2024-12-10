/** @type {import('next').NextConfig} */
const nextConfig = {
  experimental: {
      after: true
  },
  images: {
      remotePatterns: [
          {
              hostname: 'localhost'
          }
      ]
  }
};

export default nextConfig;