import type { NextConfig } from "next";

// Where the Next server proxies /api/v1 requests. Local dev defaults to the
// Go API on localhost; in Docker this is baked to the compose service
// (http://backend:8080) via the API_PROXY_TARGET build arg.
const apiTarget = process.env.API_PROXY_TARGET || "http://localhost:8080";

const nextConfig: NextConfig = {
  output: "standalone",
  rewrites: async () => {
    return {
      beforeFiles: [
        {
          source: "/api/v1/:path*",
          destination: `${apiTarget}/api/v1/:path*`,
        },
      ],
    };
  },
};

export default nextConfig;
