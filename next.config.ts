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
};

export default nextConfig;
