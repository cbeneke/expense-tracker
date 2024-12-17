# Base stage
FROM node:20-alpine AS builder

WORKDIR /app

# Install dependencies
COPY package.json package-lock.json ./
RUN npm ci

# Copy source code
COPY . .

# Development stage
FROM builder AS development
CMD ["npm", "start"]

# Build the application
FROM builder AS build
RUN npm run build

# Production stage
FROM nginx:alpine AS production

# Copy the build output
COPY --from=build /app/build /usr/share/nginx/html

# Copy nginx configuration
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Expose port
EXPOSE 80

# Start nginx
CMD ["nginx", "-g", "daemon off;"] 