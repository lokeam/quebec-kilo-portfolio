# Stage 1: Build React app
FROM node:16-alpine AS build

# Set the working directory
WORKDIR /app

# Copy package.json and package-lock.json and install dependencies
COPY package*.json ./
RUN npm install

# Copy the rest of the source code + build the app
COPY . .
RUN npm run build

# Stage 2: Serve the built static files with Nginx
FROM nginx:stable-alpine

# Remove the default Nginx static assets
RUN rm -rf /usr/share/nginx/html/*

# Copy the React build output to Nginx's public directory
COPY --from=build /app/build /usr/share/nginx/html

# Expose port 80 so that the container is accessible from the host
EXPOSE 80

# Start Nginx in the foreground
CMD ["nginx", "-g", "daemon off;"]
