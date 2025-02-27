#!/bin/bash
# Initialize Docker Swarm with advertise address
docker swarm init --advertise-addr $(hostname -I | awk '{print $1}')

# Create directory /root/data
mkdir -p /root/data

# Authenticate to Docker GHCR
docker login ghcr.io --username "${GHCR_USERNAME}" --password "${GHCR_TOKEN}"

# Pull and run the Docker container for creating api user
docker run -e DB_PATH=/app/data/database.sqlite -v /root/data/:/app/data/:consistent ghcr.io/devops-q/itu-minitwit-create-api-user:51dea8 /app/create_api_user -username=${API_USER} -password=${API_PASSWORD}


echo "Finished running minitwit init script"