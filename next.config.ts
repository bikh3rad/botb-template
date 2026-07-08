import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  /* config options here */
  output: "standalone",
  // Pin the workspace root to this project. A parent dir (/home/saeed) also
  // contains a Next.js app + lockfile, which would otherwise be inferred as
  // the root and pull in its src/proxy.ts.
  turbopack: {
    root: __dirname,
  },
  // Serve MinIO objects through the SAME ORIGIN as the site: the browser
  // requests /media/<bucket>/<key> (same host/IP it loaded the page from) and
  // the frontend proxies it to the internal MinIO endpoint. This makes media
  // URLs portable — they follow whatever server IP/host the site is accessed
  // by, instead of a baked-in localhost. Override the upstream with
  // MEDIA_UPSTREAM_URL (default the compose service address).
  async rewrites() {
    const upstream = process.env.MEDIA_UPSTREAM_URL || "http://minio:9000";
    return [
      {
        source: "/media/:path*",
        destination: `${upstream}/:path*`,
      },
    ];
  },
};

export default nextConfig;
