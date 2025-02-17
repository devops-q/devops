#!/usr/bin/env bash

HOST=$1

docker context create digitalocean --docker "host=ssh://root@152.42.150.137"

docker swarm init

docker stack deploy -c docker-compose.yml itu-minitwit