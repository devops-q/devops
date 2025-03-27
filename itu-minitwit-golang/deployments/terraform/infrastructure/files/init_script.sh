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

mkdir -p /mnt/mount

# Add permissions to the /mnt/mount directory
sudo chown -R 1000:1000 /mnt/mount

mount -o discard,defaults,noatime /dev/disk/by-id/scsi-0DO_Volume_mount /mnt/mount

echo '/dev/disk/by-id/scsi-0DO_Volume_mount /mnt/mount ext4 defaults,nofail,discard 0 0' | sudo tee -a /etc/fstab

sudo chown -R 1000:1000 /mnt/mount
sudo chmod -R 775 /mnt/mount

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


mkdir -p /root/configs

cat <<'EOF' > /root/configs/loki-config.yaml
auth_enabled: false

server:
  http_listen_port: 3100
  log_level: info

common:
  ring:
    instance_addr: 127.0.0.1
    kvstore:
      store: inmemory
  replication_factor: 1
  path_prefix: /loki

ingester:
  chunk_idle_period: 5m
  chunk_target_size: 1048576
  max_chunk_age: 1h
  lifecycler:
    ring:
      kvstore:
        store: inmemory
      replication_factor: 1
    heartbeat_period: 1m
  wal:
    enabled: true
    dir: /loki/wal

frontend:
  address: 127.0.0.1
  max_outstanding_per_tenant: 2048
  compress_responses: true

frontend_worker:
  frontend_address: 127.0.0.1
  grpc_client_config:
    max_send_msg_size: 104857600
    max_recv_msg_size: 104857600

query_scheduler:
  max_outstanding_requests_per_tenant: 4096
  use_scheduler_ring: true

schema_config:
  configs:
    - from: 2020-10-24
      store: tsdb
      object_store: s3
      schema: v13
      index:
        prefix: index_
        period: 24h

storage_config:
  tsdb_shipper:
    active_index_directory: /loki/index
    cache_location: /loki/index_cache
  aws:
    endpoint: fra1.digitaloceanspaces.com
    access_key_id: '${S3_ACCESS_KEY}'
    secret_access_key: '${S3_SECRET_KEY}'
    s3forcepathstyle: true
    insecure: false
    region: us-east-1
    bucketnames: '${S3_BUCKET_NAME}'

compactor:
  working_directory: /loki/compactor
  shared_store: s3
  compaction_interval: 10m
EOF

cat <<'EOF' > /root/configs/alloy-config.alloy
discovery.docker "linux" {
  host = "unix:///var/run/docker.sock"
}

discovery.relabel "container_labels" {
  targets = discovery.docker.linux.targets

  rule {
    source_labels = ["__meta_docker_container_name"]
    target_label  = "container_name"
  }

  rule {
    source_labels = ["__meta_docker_container_id"]
    target_label  = "container_id"
  }

  rule {
    source_labels = ["__meta_docker_container_label_com_docker_swarm_service_name"]
    target_label  = "service_name"
  }
}

loki.source.docker "default" {
  host          = "unix:///var/run/docker.sock"
  targets       = discovery.relabel.container_labels.output
  labels        = {"app" = "docker"}
  forward_to    = [loki.write.local.receiver]
  refresh_interval = "5s"
}

loki.write "local" {
  endpoint {
    url = "http://loki:3100/loki/api/v1/push"
  }
}
EOF


echo "Finished running minitwit init script"
