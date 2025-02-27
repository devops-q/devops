#!/bin/bash
# Initialize Docker Swarm with advertise address
docker swarm init --advertise-addr $(hostname -I | awk '{print $1}')

# Create directory /root/data
mkdir -p /root/data

# Authenticate to Docker GHCR
echo "${GHCR_TOKEN}" | docker login ghcr.io -u "${GHCR_USERNAME}" --password-stdin

# Pull and run the Docker container for creating api user
docker pull ghcr.io/${GHCR_USERNAME}/your-container:latest

docker run -d -e DB_PATH=/root/data/database.sqlite ghcr.io/devops-q/itu-minitwit-create-api-user:51dea8 /app/create_api_user -username=${API_USER} -password=${API_PASSWORD}