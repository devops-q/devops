#!/bin/bash
# Initialize Docker Swarm with advertise address
docker swarm init --advertise-addr $(hostname -I | awk '{print $1}')


# Allow ssh traffic on port 22
sudo ufw allow 22/tcp
sudo ufw allow 22

# Create directory /root/data
mkdir -p /root/data

# Pull and run the Docker container for creating api user
docker run -e DB_HOST=${DB_HOST} \
 -e DB_USER=${DB_USER} \
 -e DB_PASSWORD=${DB_PASSWORD} \
 -e DB_NAME=${DB_NAME} \
 -e DB_PORT=${DB_PORT} \
 -e DB_SSL_MODE=require \
 -e API_USER=${API_USER} \
 -e API_PASSWORD=${API_PASSWORD} \
 ghcr.io/devops-q/itu-minitwit-create-api-user:dec2f8


echo "Finished running minitwit init script"