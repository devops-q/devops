variable "vm_ip_address" {
  type = string
}

variable "grafana_root_password" {
  type = string
  sensitive = true
}

terraform {
  required_providers {
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
    key = "tf-observability.tfstate"

    # Deactivate a few AWS-specific checks
    skip_credentials_validation = true
    skip_requesting_account_id  = true
    skip_metadata_api_check     = true
    skip_region_validation      = true
    skip_s3_checksum            = true
    region                      = "us-east-1"
  }

}


provider "grafana" {
  url  = "http://${var.vm_ip_address}:3000"
  auth = "admin:${var.grafana_root_password}"
}

