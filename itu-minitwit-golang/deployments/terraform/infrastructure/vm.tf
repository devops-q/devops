variable "prometheus_root_password" {
  type        = string
  description = "API password for the initial user"
}

variable "helge_and_mircea_password" {
  type        = string
  description = "API password for the initial user"
}

variable "s3_logs_bucket_name" {
  type        = string
  description = "S3 bucket name for logs"
}

variable "s3_access_key" {
  type        = string
  description = "S3 access key"
}

variable "s3_secret_key" {
  type        = string
  description = "S3 secret key"
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
    S3_ACCESS_KEY = var.s3_access_key
    S3_SECRET_KEY = var.s3_secret_key
    S3_BUCKET_NAME = var.s3_logs_bucket_name
  })

  lifecycle {
    create_before_destroy = true
    ignore_changes = [user_data] # Ignore changes to user_data so we don't have to recreate the droplet.
  }
}

resource "digitalocean_floating_ip" "ip" {
  droplet_id = digitalocean_droplet.minitwit-vm.id
  region     = digitalocean_droplet.minitwit-vm.region
}

output "ip_address" {
  value = trimspace(digitalocean_floating_ip.ip.ip_address)
}
