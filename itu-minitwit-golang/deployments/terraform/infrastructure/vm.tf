resource "digitalocean_droplet" "minitwit-vm" {
  image  = "docker-20-04"
  name   = "minitwit-vm"
  region = "fra1"
  size   = "s-1vcpu-2gb-70gb-intel"
  ssh_keys = [
    data.digitalocean_ssh_key.terraform.id
  ]
  lifecycle {
    create_before_destroy = true
    ignore_changes = [user_data] # Ignore changes to user_data so we don't have to recreate the droplet.
  }
}

resource "digitalocean_volume" "mount" {
  region = "fra1"
  name = "mount"
  size = "20"
  initial_filesystem_type = "ext4"

}

resource "digitalocean_volume_attachment" "mount" {
  droplet_id = digitalocean_droplet.minitwit-vm.id
  volume_id = digitalocean_volume.mount.id

}



resource "digitalocean_floating_ip" "ip" {
  droplet_id = digitalocean_droplet.minitwit-vm.id
  region     = digitalocean_droplet.minitwit-vm.region
}

output "ip_address" {
  value = trimspace(digitalocean_floating_ip.ip.ip_address)
}
