import { hostname } from "os";

/** @type {import('next').NextConfig} */
const nextConfig = {
  experimental: {
      after: true
  },
  images: {
      remotePatterns: [
          {
              hostname: 'localhost'
          },
          {
            hostname: 'host.docker.internal'
          }
      ]
  }
};

export default nextConfig;