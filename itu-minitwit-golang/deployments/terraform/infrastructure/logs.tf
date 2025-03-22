variable "s3_access_key" {
  type        = string
  description = "S3 access key"
}

variable "s3_secret_key" {
  type        = string
  description = "S3 secret key"
}

resource "random_id" "server" {
  byte_length = 8
}

resource "digitalocean_spaces_bucket" "logs" {
  name   = "loki-logs-${random_id.server.hex}"
  region = "fra1"
}

