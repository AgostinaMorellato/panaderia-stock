# Use the official Node image
FROM node:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the local package files to the container's workspace
COPY . .

ENV PUPPETEER_SKIP_DOWNLOAD true

# Install dependencies
RUN npm install

# Expose port 3000 to the outside world
EXPOSE 3001

# Command to run the application
CMD ["npm", "start"]