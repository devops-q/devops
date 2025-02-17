resource "digitalocean_droplet" "minitwit-vm" {
  image  = "docker-20-04"
  name   = "minitwit-vm"
  region = "ams3"
  size   = "s-2vcpu-4gb-120gb-intel"
  ssh_keys = [
    data.digitalocean_ssh_key.terraform.id
  ]
}

resource "digitalocean_floating_ip" "ip" {
  droplet_id = digitalocean_droplet.minitwit-vm.id
  region     = digitalocean_droplet.minitwit-vm.region
}