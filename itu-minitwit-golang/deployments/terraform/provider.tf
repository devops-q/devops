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
  url  = "http://${data.digitalocean_reserved_ip.terraform.id}:3000"
  auth = "admin:admin"
}


// Create resources (optional: within the organization)
resource "grafana_folder" "my_folder" {
  title  = "grafana_dashboard_folder"
}

resource "grafana_dashboard" "test_folder" {
  folder = grafana_folder.my_folder.id
  config_json = jsonencode({
    "annotations": {
      "list": [
        {
          "builtIn": 1,
          "datasource": {
            "type": "grafana",
            "uid": "-- Grafana --"
          },
          "enable": true,
          "hide": true,
          "iconColor": "rgba(0, 211, 255, 1)",
          "name": "Annotations & Alerts",
          "type": "dashboard"
        }
      ]
    },
    "editable": true,
    "fiscalYearStartMonth": 0,
    "graphTooltip": 0,
    "id": 3,
    "links": [],
    "liveNow": false,
    "panels": [
      {
        "datasource": {
          "type": "prometheus",
          "uid": "b957054d-2c8d-403e-95ef-49a1aea7f57b"
        },
        "fieldConfig": {
          "defaults": {
            "color": {
              "mode": "palette-classic"
            },
            "custom": {
              "axisBorderShow": false,
              "axisCenteredZero": false,
              "axisColorMode": "text",
              "axisLabel": "",
              "axisPlacement": "auto",
              "barAlignment": 0,
              "drawStyle": "line",
              "fillOpacity": 0,
              "gradientMode": "none",
              "hideFrom": {
                "legend": false,
                "tooltip": false,
                "viz": false
              },
              "insertNulls": false,
              "lineInterpolation": "linear",
              "lineWidth": 1,
              "pointSize": 5,
              "scaleDistribution": {
                "type": "linear"
              },
              "showPoints": "auto",
              "spanNulls": false,
              "stacking": {
                "group": "A",
                "mode": "none"
              },
              "thresholdsStyle": {
                "mode": "off"
              }
            },
            "mappings": [],
            "thresholds": {
              "mode": "absolute",
              "steps": [
                {
                  "color": "green",
                  "value": null
                },
                {
                  "color": "red",
                  "value": 80
                }
              ]
            }
          },
          "overrides": []
        },
        "gridPos": {
          "h": 8,
          "w": 12,
          "x": 0,
          "y": 0
        },
        "id": 2,
        "options": {
          "legend": {
            "calcs": [],
            "displayMode": "list",
            "placement": "bottom",
            "showLegend": true
          },
          "tooltip": {
            "mode": "single",
            "sort": "none"
          }
        },
        "targets": [
          {
            "datasource": {
              "type": "prometheus",
              "uid": "Prometheus"
            },
            "disableTextWrap": false,
            "editorMode": "builder",
            "expr": "sum(rate(gin_gin_requests_total[$__rate_interval]))",
            "fullMetaSearch": false,
            "includeNullMetadata": false,
            "instant": false,
            "legendFormat": "__auto",
            "range": true,
            "refId": "A",
            "useBackend": false
          }
        ],
        "title": "Rate of requests",
        "type": "timeseries"
      },
      {
        "datasource": {
          "type": "prometheus",
          "uid": "Prometheus"
        },
        "fieldConfig": {
          "defaults": {
            "color": {
              "mode": "thresholds"
            },
            "mappings": [],
            "thresholds": {
              "mode": "absolute",
              "steps": [
                {
                  "color": "green",
                  "value": null
                },
                {
                  "color": "red",
                  "value": 80
                }
              ]
            }
          },
          "overrides": []
        },
        "gridPos": {
          "h": 8,
          "w": 12,
          "x": 12,
          "y": 0
        },
        "id": 1,
        "options": {
          "colorMode": "value",
          "graphMode": "area",
          "justifyMode": "auto",
          "orientation": "auto",
          "reduceOptions": {
            "calcs": [
              "lastNotNull"
            ],
            "fields": "",
            "values": false
          },
          "textMode": "auto",
          "wideLayout": true
        },
        "pluginVersion": "10.2.4",
        "targets": [
          {
            "datasource": {
              "type": "prometheus",
              "uid": "Prometheus"
            },
            "disableTextWrap": false,
            "editorMode": "builder",
            "expr": "sum(gin_gin_requests_total)",
            "fullMetaSearch": false,
            "includeNullMetadata": false,
            "instant": false,
            "legendFormat": "__auto",
            "range": true,
            "refId": "A",
            "useBackend": false
          }
        ],
        "title": "Number of requests on page",
        "type": "stat"
      }
    ],
    "refresh": "",
    "schemaVersion": 39,
    "tags": [],
    "templating": {
      "list": []
    },
    "time": {
      "from": "now-6h",
      "to": "now"
    },
    "timepicker": {},
    "timezone": "",
    "title": "Main dashboard",
    "uid": "e115a275-682c-4ec2-8482-55e552c2c3a0",
    "version": 2,
    "weekStart": ""
  })
}