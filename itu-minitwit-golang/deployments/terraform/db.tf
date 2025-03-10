variable "db_name" {
  type        = string
  description = "The name of the PostgreSQL database"
}

resource "digitalocean_database_cluster" "postgres" {
  name       = "minitwit-db"
  engine     = "pg"
  version    = "17"
  size       = "db-s-1vcpu-2gb"
  region     = "fra1"
  node_count = 1
}

resource "digitalocean_database_firewall" "minitwit_app_firewall" {
  cluster_id = digitalocean_database_cluster.postgres.id

  rule {
    type  = "droplet"
    value = digitalocean_droplet.minitwit-vm.id
  }
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
  value = digitalocean_database_cluster.postgres.user
}

output "db_password" {
  value     = digitalocean_database_cluster.postgres.password
  sensitive = true
}