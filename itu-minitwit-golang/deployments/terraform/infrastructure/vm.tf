variable "prometheus_root_password" {
  type        = string
  description = "API password for the initial user"
}

variable "helge_and_mircea_password" {
  type        = string
  description = "API password for the initial user"
}

resource "digitalocean_droplet" "minitwit-vm" {
  image  = "docker-20-04"
  name   = "minitwit-vm"
  region = "fra1"
  size   = "s-1vcpu-2gb-70gb-intel"
  ssh_keys = [
    data.digitalocean_ssh_key.terraform.id
  ]
  user_data = templatefile("./files/init_script.sh", {
    PROMETHEUS_ROOT_PASSWORD = var.prometheus_root_password
    PROMETHEUS_ROOT_PASSWORD_BCRYPT = bcrypt(var.prometheus_root_password) # Prometheus expects a bcrypt hash
    HELGE_AND_MIRCEA_PASSWORD_BCRYPT = bcrypt(var.helge_and_mircea_password) # Prometheus expects a bcrypt hash
  })

  lifecycle {
    create_before_destroy = true
  }
}

resource "digitalocean_floating_ip" "ip" {
  droplet_id = digitalocean_droplet.minitwit-vm.id
  region     = digitalocean_droplet.minitwit-vm.region
}

output "ip_address" {
  value = trimspace(digitalocean_floating_ip.ip.ip_address)
}
