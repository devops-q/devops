variable "db_name" {
  type        = string
  description = "The name of the PostgreSQL database"
}

variable "db_user" {
  type        = string
  description = "The username for the PostgreSQL database"
}

variable "db_password" {
  type        = string
  description = "The password for the PostgreSQL database"
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
  password   = var.db_password
}

resource "digitalocean_database_db" "app_db" {
  cluster_id = digitalocean_database_cluster.postgres.id
  name       = var.db_name
}

output "db_host" {
  value = digitalocean_database_cluster.postgres.host
}

output "db_port" {
  value = digitalocean_database_cluster.postgres.port
}