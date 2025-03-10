#!/usr/bin/env bash

timeout=360
elapsed=0
while ! grep -q "Finished running minitwit init script" /var/log/cloud-init-output.log; do
  if [ $elapsed -ge $timeout ]; then
    echo "Timeout reached. Exiting."
    exit 1
  fi
  echo "Waiting for user_data script to finish..."
  sleep 10
  elapsed=$((elapsed + 10))
done