#!/bin/bash
# Initialize Docker Swarm with advertise address
docker swarm init --advertise-addr $(hostname -I | awk '{print $1}')


# Allow ssh traffic on port 22
sudo ufw allow 22/tcp
sudo ufw allow 22

# Create directory /root/data
mkdir -p /root/data

# Create and populate the
mkdir -p /root/prometheus

touch /root/prometheus

echo "global:
        scrape_interval:     15s # By default, scrape targets every 15 seconds.
        evaluation_interval: 15s # Evaluate rules every 15 seconds.

        # Attach these extra labels to all timeseries collected by this Prometheus instance.
        external_labels:
          monitor: 'codelab-monitor'

      rule_files:
        - 'prometheus.rules.yml'

      scrape_configs:
        - job_name: 'prometheus'

          # Override the global default and scrape targets from this job every 5 seconds.
          scrape_interval: 5s

          static_configs:
            - targets: ['prometheus:9090']

        - job_name:       'itu-minitwit-app'

          # Override the global default and scrape targets from this job every 5 seconds.
          scrape_interval: 5s

          static_configs:
            - targets: ['app:80']
              labels:
                group: 'production'" > /root/prometheus

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