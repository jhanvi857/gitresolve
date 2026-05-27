/** @type {import('next').NextConfig} */
const nextConfig = {
  async redirects() {
    return [
      {
        source: '/docs',
        destination: '/docs/installation',
        permanent: true,
      },
    ]
  },
};

export default nextConfig;
