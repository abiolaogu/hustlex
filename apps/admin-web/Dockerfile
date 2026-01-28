# =============================================================================
# HustleX Admin Dashboard Dockerfile
# =============================================================================
# Multi-stage build for React + Refine v4 admin dashboard
# Build: docker build -t hustlex-admin .
# Run: docker run -p 3000:3000 hustlex-admin

# -----------------------------------------------------------------------------
# Stage 1: Dependencies
# -----------------------------------------------------------------------------
FROM node:20-alpine AS deps

WORKDIR /app

# Copy package files
COPY package.json package-lock.json* ./

# Install dependencies (use ci if lock file exists, otherwise install)
RUN --mount=type=cache,target=/root/.npm \
    if [ -f package-lock.json ]; then npm ci --prefer-offline; else npm install; fi

# -----------------------------------------------------------------------------
# Stage 2: Builder
# -----------------------------------------------------------------------------
FROM node:20-alpine AS builder

WORKDIR /app

# Copy dependencies from deps stage
COPY --from=deps /app/node_modules ./node_modules

# Copy source files
COPY . .

# Build arguments for environment
ARG VITE_API_URL=http://localhost:8081
ARG VITE_GRAPHQL_URL=http://localhost:8080/v1/graphql
ARG VITE_HASURA_ADMIN_SECRET=hustlex_hasura_admin_secret
ARG VITE_WS_URL=ws://localhost:8080/v1/graphql

# Set environment for build
ENV VITE_API_URL=$VITE_API_URL
ENV VITE_GRAPHQL_URL=$VITE_GRAPHQL_URL
ENV VITE_HASURA_ADMIN_SECRET=$VITE_HASURA_ADMIN_SECRET
ENV VITE_WS_URL=$VITE_WS_URL

# Build the application
RUN npm run build

# -----------------------------------------------------------------------------
# Stage 3: Production
# -----------------------------------------------------------------------------
FROM nginx:1.25-alpine AS production

# Install envsubst for runtime environment variable injection
RUN apk add --no-cache gettext

# Create non-root user
RUN addgroup -g 1000 -S hustlex && \
    adduser -u 1000 -S hustlex -G hustlex

# Copy built assets
COPY --from=builder /app/dist /usr/share/nginx/html

# Copy nginx configuration
COPY nginx.conf /etc/nginx/nginx.conf

# Create env config script for runtime environment variables
RUN cat > /docker-entrypoint.d/40-env-config.sh << 'EOF'
#!/bin/sh
set -e

# Generate runtime config from environment variables
cat > /usr/share/nginx/html/config.js << JSEOF
window.__ENV__ = {
  VITE_API_URL: "${VITE_API_URL:-http://localhost:8081}",
  VITE_GRAPHQL_URL: "${VITE_GRAPHQL_URL:-http://localhost:8080/v1/graphql}",
  VITE_WS_URL: "${VITE_WS_URL:-ws://localhost:8080/v1/graphql}",
  VITE_HASURA_ADMIN_SECRET: "${VITE_HASURA_ADMIN_SECRET:-}"
};
JSEOF

echo "Runtime config generated"
EOF

RUN chmod +x /docker-entrypoint.d/40-env-config.sh

# Set permissions
RUN chown -R hustlex:hustlex /usr/share/nginx/html && \
    chown -R hustlex:hustlex /var/cache/nginx && \
    chown -R hustlex:hustlex /var/log/nginx && \
    touch /var/run/nginx.pid && \
    chown -R hustlex:hustlex /var/run/nginx.pid

# Switch to non-root user (nginx needs root for port binding in alpine)
# USER hustlex

# Expose port
EXPOSE 3000

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget -q --spider http://localhost:3000/health || exit 1

# Start nginx
CMD ["nginx", "-g", "daemon off;"]

# -----------------------------------------------------------------------------
# Stage 4: Development
# -----------------------------------------------------------------------------
FROM node:20-alpine AS development

WORKDIR /app

# Install dependencies
COPY package.json package-lock.json* ./
RUN npm ci

# Copy source (will be overridden by volume mount)
COPY . .

# Expose Vite dev server port
EXPOSE 3000

# Start development server
CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0"]
