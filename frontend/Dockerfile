# syntax=docker/dockerfile:1

# ---- deps ----
FROM node:22-alpine AS deps
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci

# ---- build ----
FROM node:22-alpine AS build
WORKDIR /app
# Baked into the route manifest at build time; points the /api/v1 proxy at the
# backend service on the compose network.
ARG API_PROXY_TARGET=http://localhost:8080
ENV API_PROXY_TARGET=$API_PROXY_TARGET
COPY --from=deps /app/node_modules ./node_modules
COPY . .
RUN npm run build

# ---- dev (hot reload; used by docker-compose.override.yml) ----
# Source is bind-mounted and node_modules live in a named volume (installed at
# container start), so edits are picked up live with no image rebuild.
FROM node:22-alpine AS dev
WORKDIR /app
ENV NODE_ENV=development NEXT_TELEMETRY_DISABLED=1
EXPOSE 3000
CMD ["npm", "run", "dev", "--", "-H", "0.0.0.0"]

# ---- runtime (Next standalone) ----
FROM node:22-alpine AS runner
WORKDIR /app
ENV NODE_ENV=production PORT=3000 HOSTNAME=0.0.0.0
# NOTE: the /api/v1 proxy target (API_PROXY_TARGET) is consumed at BUILD time by
# next.config.ts — Next bakes rewrite destinations into the route manifest, so it
# is NOT runtime-configurable. Set it as a build arg/variable pointing at the
# backend service reachable from the proxy on the deployment network.
RUN addgroup -S nodejs && adduser -S nextjs -G nodejs
COPY --from=build /app/public ./public
COPY --from=build /app/.next/standalone ./
COPY --from=build /app/.next/static ./.next/static
USER nextjs
EXPOSE 3000
CMD ["node", "server.js"]
