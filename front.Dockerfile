# Stage 1: Build the Vue app
FROM node:18-alpine AS build

# Set working directory inside container
WORKDIR /app

# Copy package files first (for caching installs)
COPY frontend/package*.json ./
RUN npm install

# Copy the rest of your frontend source
COPY frontend/ ./

# Build the app
RUN npm run build

# Stage 2: Serve with a lightweight static server
FROM node:18-alpine

WORKDIR /app

# Install a simple static file server
RUN npm install -g serve

# Copy built frontend from build stage
COPY --from=build /app/dist ./dist

# Expose port (Traefik will route traffic here)
EXPOSE 3000

# Start the server
CMD ["serve", "-s", "dist", "-l", "3000"]
