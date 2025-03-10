terraform {
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "2.48.2"
    }
    grafana = {
      source  = "grafana/grafana"
      version = "3.15.2"
    }
  }

  backend "s3" {
    endpoints = {
      s3 = "https://fra1.digitaloceanspaces.com"
    }
    bucket = "itu-devops-q-remote-prod"
    key = "tf.tfstate"

    # Deactivate a few AWS-specific checks
    skip_credentials_validation = true
    skip_requesting_account_id  = true
    skip_metadata_api_check     = true
    skip_region_validation      = true
    skip_s3_checksum            = true
    region                      = "us-east-1"
  }

}

variable "do_token" {}
variable "do_ssh_key_name" {}



provider "digitalocean" {
  token = var.do_token
}

data "digitalocean_ssh_key" "terraform" {
  name = var.do_ssh_key_name
}


provider "grafana" {
  url  = "http://164.90.242.193:3000"
  auth = "admin:admin"
}


