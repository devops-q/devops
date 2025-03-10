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

variable "do_float_ip" {}

provider "digitalocean" {
  token = var.do_token
}

data "digitalocean_ssh_key" "terraform" {
  name = var.do_ssh_key_name
}


data "digitalocean_reserved_ip" "terraform" {
  ip_address = var.do_float_ip
}

provider "grafana" {
  url  = "http://${var.do_float_ip}:3000"
  auth = "admin:admin"
}

// Optional (On-premise, not supported in Grafana Cloud): Create an organization
resource "grafana_organization" "my_org" {
  name = "my_org"
}

// Create resources (optional: within the organization)
resource "grafana_folder" "my_folder" {
  org_id = grafana_organization.my_org.org_id
  title  = "Test Folder"
}

resource "grafana_dashboard" "test_folder" {
  org_id = grafana_organization.my_org.org_id
  folder = grafana_folder.my_folder.id
  config_json = jsonencode({
    "title" : "My Dashboard Title",
    "uid" : "my-dashboard-uid"
    // ... other dashboard properties
  })
}