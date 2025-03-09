variable "db_name" {
  type        = string
  description = "The name of the PostgreSQL database"
}

variable "db_user" {
  type        = string
  description = "The username for the PostgreSQL database"
}


resource "digitalocean_database_cluster" "postgres" {
  name       = "minitwit-db"
  engine     = "pg"
  version    = "17"
  size       = "db-s-1vcpu-2gb"
  region     = "fra1"
  node_count = 1
}

resource "digitalocean_database_user" "app_user" {
  cluster_id = digitalocean_database_cluster.postgres.id
  name       = var.db_user
}

resource "digitalocean_database_db" "app_db" {
  cluster_id = digitalocean_database_cluster.postgres.id
  name       = var.db_name
}

output "db_private_host" {
  value = digitalocean_database_cluster.postgres.private_host
}

output "db_port" {
  value = digitalocean_database_cluster.postgres.port
}

output "db_user" {
  value = digitalocean_database_user.app_user.name
}

output "db_password" {
  value = digitalocean_database_user.app_user.password
  sensitive = true
}