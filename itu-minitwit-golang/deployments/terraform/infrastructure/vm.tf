variable "api_user" {
  type        = string
  description = "Initial API user to be created"
}

variable "api_password" {
  type        = string
  description = "API password for the initial user"
}

resource "grafana_data_source" "prometheus" {
  type                = "prometheus"
  name                = "mimir"
  url                 = "https://prometheus:9090"
  basic_auth_enabled  = true
  basic_auth_username = "admin"

  json_data_encoded = jsonencode({
    httpMethod        = "POST"
    prometheusType    = "Mimir"
    prometheusVersion = "2.4.0"
  })

  uid = "prometheus"

  secure_json_data_encoded = jsonencode({
    basicAuthPassword = "admin"
  })
  depends_on = [digitalocean_droplet.minitwit-vm]
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
    API_USER     = var.api_user
    API_PASSWORD = var.api_password
    DB_HOST      = digitalocean_database_cluster.postgres.private_host
    DB_USER      = digitalocean_database_cluster.postgres.user
    DB_PASSWORD  = digitalocean_database_cluster.postgres.password
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




resource "grafana_folder" "my_folder" {
  title  = "grafana_dashboard_folder"

  depends_on = [digitalocean_droplet.minitwit-vm]
}



resource "grafana_dashboard" "grafana_dashboard_folder" {
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
          "uid": grafana_data_source.prometheus.uid
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
              "uid": grafana_data_source.prometheus.uid
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
          "uid": "prometheus"
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
              "uid": grafana_data_source.prometheus.uid
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
  depends_on = [grafana_folder.my_folder, digitalocean_droplet.minitwit-vm]

}

