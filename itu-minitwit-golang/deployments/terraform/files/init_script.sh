#!/bin/bash
# Initialize Docker Swarm with advertise address
docker swarm init --advertise-addr $(hostname -I | awk '{print $1}')

# Create directory /root/data
mkdir -p /root/data

# Pull and run the Docker container for creating api user
docker run -e DB_HOST=${DB_HOST} \
 -e DB_USER=${DB_USER} \
 -e DB_PASSWORD=${DB_PASSWORD} \
 -e DB_NAME=${DB_NAME} \
 -e DB_PORT=${DB_PORT} \
 ghcr.io/devops-q/itu-minitwit-create-api-user:c1960e /app/create_api_user \
 -username=${API_USER} \
 -password=${API_PASSWORD}


echo "Finished running minitwit init script"