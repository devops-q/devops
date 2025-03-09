variable "api_user" {
  type        = string
  description = "Initial API user to be created"
}

variable "api_password" {
  type        = string
  description = "API password for the initial user"
}

resource "digitalocean_droplet" "minitwit-vm" {
  image  = "docker-20-04"
  name   = "minitwit-vm"
  region = "fra1"
  size   = "s-2vcpu-4gb-120gb-intel"
  ssh_keys = [
    data.digitalocean_ssh_key.terraform.id
  ]
  user_data = templatefile("./files/init_script.sh", {
    API_USER     = var.api_user
    API_PASSWORD = var.api_password
    DB_HOST      = digitalocean_database_cluster.postgres.private_host
    DB_USER      = digitalocean_database_user.app_user.name
    DB_PASSWORD  = digitalocean_database_user.app_user.password
    DB_NAME      = digitalocean_database_db.app_db.name
    DB_PORT      = digitalocean_database_cluster.postgres.port
  })
}

resource "digitalocean_floating_ip" "ip" {
  droplet_id = digitalocean_droplet.minitwit-vm.id
  region     = digitalocean_droplet.minitwit-vm.region
}

output "ip_address" {
  value = trimspace(digitalocean_floating_ip.ip.ip_address)
}