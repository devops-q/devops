#!/bin/bash
# Initialize Docker Swarm with advertise address
docker swarm init --advertise-addr "$(hostname -I | awk '{print $1}')"


# Allow ssh traffic on port 22
sudo ufw allow 22/tcp
sudo ufw allow 22

# Create directory /root/data
mkdir -p /root/data

# Create and populate the
mkdir -p /root/prometheus

cat <<'EOF' > /root/prometheus/prometheus.yml
global:
  scrape_interval: 15s  # By default, scrape targets every 15 seconds.
  evaluation_interval: 15s  # Evaluate rules every 15 seconds.
  # Attach these extra labels to all timeseries collected by this Prometheus instance.
  external_labels:
    monitor: 'codelab-monitor'

rule_files:
  - 'prometheus.rules.yml'

scrape_configs:
  - job_name: 'prometheus'
    scrape_interval: 5s  # Scrape targets every 5 seconds for this job.
    static_configs:
      - targets: ['prometheus:9090']
    basic_auth:
      username: 'admin'
      password: '${PROMETHEUS_ROOT_PASSWORD}'

  - job_name: 'itu-minitwit-app'
    scrape_interval: 5s  # Scrape targets every 5 seconds for this job.
    static_configs:
      - targets: ['app:80']
        labels:
          group: 'production'
EOF

cat <<'EOF' > /root/prometheus/web.yml
basic_auth_users:
    admin: '${PROMETHEUS_ROOT_PASSWORD_BCRYPT}'
    helgeandmircea: '${HELGE_AND_MIRCEA_PASSWORD_BCRYPT}'
EOF

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
