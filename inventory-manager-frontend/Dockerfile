# Build stage
FROM node:20-alpine AS build

WORKDIR /app

# Build-time env vars for Next.js (can be provided via --build-arg)
ARG NEXT_PUBLIC_API_BASE_URL
ARG NEXT_PUBLIC_WS_URL

# Expose build args as env vars so Next.js can read process.env.NEXT_PUBLIC_*
ENV NEXT_PUBLIC_API_BASE_URL=$NEXT_PUBLIC_API_BASE_URL
ENV NEXT_PUBLIC_WS_URL=$NEXT_PUBLIC_WS_URL

# Copy package files
COPY package.json package-lock.json ./

# Install dependencies
RUN npm ci

# Copy source code
COPY . .

# Build the Next.js application
RUN npm run build

# Production stage
FROM node:20-alpine AS production

WORKDIR /app

# Set environment to production
ENV NODE_ENV=production

# Copy the standalone output from build stage
COPY --from=build /app/.next/standalone ./
COPY --from=build /app/.next/static ./.next/static
COPY --from=build /app/public ./public

# Expose port
EXPOSE 3000

# Start the application (listen on all interfaces)
# Note: HOSTNAME is set at runtime to override Kubernetes auto-injection
ENV PORT=3000
CMD ["sh", "-c", "HOSTNAME=0.0.0.0 node server.js"]

