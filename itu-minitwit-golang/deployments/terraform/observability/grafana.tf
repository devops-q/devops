variable "helge_and_mircea_password" {
  type      = string
  sensitive = true
}

variable "prometheus_root_password" {
  type      = string
  sensitive = true
}


resource "grafana_folder" "my_folder" {
  title = "grafana_dashboard_folder"
}

resource "grafana_data_source" "loki" {
  name = "loki"
  type = "loki"
  url = "http://loki:3100"

  uid = "loki"
}

resource "grafana_data_source" "prometheus" {
  type                = "prometheus"
  name                = "mimir"
  url                 = "http://prometheus:9090"
  basic_auth_enabled  = true
  basic_auth_username = "admin"

  json_data_encoded = jsonencode({
    httpMethod        = "POST"
    prometheusType    = "Mimir"
    prometheusVersion = "2.4.0"
  })

  uid = "prometheus"

  secure_json_data_encoded = jsonencode({
    basicAuthPassword = var.prometheus_root_password
  })
}


resource "grafana_dashboard" "grafana_dashboard_folder" {
  folder = grafana_folder.my_folder.id
  config_json = jsonencode({
    "annotations" : {
      "list" : [
        {
          "builtIn" : 1,
          "datasource" : {
            "type" : "grafana",
            "uid" : "-- Grafana --"
          },
          "enable" : true,
          "hide" : true,
          "iconColor" : "rgba(0, 211, 255, 1)",
          "name" : "Annotations & Alerts",
          "type" : "dashboard"
        }
      ]
    },
    "editable" : true,
    "fiscalYearStartMonth" : 0,
    "graphTooltip" : 0,
    "id" : 3,
    "links" : [],
    "liveNow" : false,
    "panels" : [
      {
        "datasource" : {
          "type" : "prometheus",
          "uid" : grafana_data_source.prometheus.uid
        },
        "fieldConfig" : {
          "defaults" : {
            "color" : {
              "mode" : "palette-classic"
            },
            "custom" : {
              "axisBorderShow" : false,
              "axisCenteredZero" : false,
              "axisColorMode" : "text",
              "axisLabel" : "",
              "axisPlacement" : "auto",
              "barAlignment" : 0,
              "drawStyle" : "line",
              "fillOpacity" : 0,
              "gradientMode" : "none",
              "hideFrom" : {
                "legend" : false,
                "tooltip" : false,
                "viz" : false
              },
              "insertNulls" : false,
              "lineInterpolation" : "linear",
              "lineWidth" : 1,
              "pointSize" : 5,
              "scaleDistribution" : {
                "type" : "linear"
              },
              "showPoints" : "auto",
              "spanNulls" : false,
              "stacking" : {
                "group" : "A",
                "mode" : "none"
              },
              "thresholdsStyle" : {
                "mode" : "off"
              }
            },
            "mappings" : [],
            "thresholds" : {
              "mode" : "absolute",
              "steps" : [
                {
                  "color" : "green",
                  "value" : null
                },
                {
                  "color" : "red",
                  "value" : 80
                }
              ]
            }
          },
          "overrides" : []
        },
        "gridPos" : {
          "h" : 8,
          "w" : 12,
          "x" : 0,
          "y" : 0
        },
        "id" : 2,
        "options" : {
          "legend" : {
            "calcs" : [],
            "displayMode" : "list",
            "placement" : "bottom",
            "showLegend" : true
          },
          "tooltip" : {
            "mode" : "single",
            "sort" : "none"
          }
        },
        "targets" : [
          {
            "datasource" : {
              "type" : "prometheus",
              "uid" : grafana_data_source.prometheus.uid
            },
            "disableTextWrap" : false,
            "editorMode" : "builder",
            "expr" : "sum(rate(gin_gin_requests_total[$__rate_interval]))",
            "fullMetaSearch" : false,
            "includeNullMetadata" : false,
            "instant" : false,
            "legendFormat" : "__auto",
            "range" : true,
            "refId" : "A",
            "useBackend" : false
          }
        ],
        "title" : "Rate of requests",
        "type" : "timeseries"
      },
      {
        "datasource" : {
          "type" : "prometheus",
          "uid" : "prometheus"
        },
        "fieldConfig" : {
          "defaults" : {
            "color" : {
              "mode" : "thresholds"
            },
            "mappings" : [],
            "thresholds" : {
              "mode" : "absolute",
              "steps" : [
                {
                  "color" : "green",
                  "value" : null
                },
                {
                  "color" : "red",
                  "value" : 80
                }
              ]
            }
          },
          "overrides" : []
        },
        "gridPos" : {
          "h" : 8,
          "w" : 12,
          "x" : 12,
          "y" : 0
        },
        "id" : 1,
        "options" : {
          "colorMode" : "value",
          "graphMode" : "area",
          "justifyMode" : "auto",
          "orientation" : "auto",
          "reduceOptions" : {
            "calcs" : [
              "lastNotNull"
            ],
            "fields" : "",
            "values" : false
          },
          "textMode" : "auto",
          "wideLayout" : true
        },
        "pluginVersion" : "10.2.4",
        "targets" : [
          {
            "datasource" : {
              "type" : "prometheus",
              "uid" : grafana_data_source.prometheus.uid
            },
            "disableTextWrap" : false,
            "editorMode" : "builder",
            "expr" : "sum(gin_gin_requests_total)",
            "fullMetaSearch" : false,
            "includeNullMetadata" : false,
            "instant" : false,
            "legendFormat" : "__auto",
            "range" : true,
            "refId" : "A",
            "useBackend" : false
          }
        ],
        "title" : "Number of requests on page",
        "type" : "stat"
      }
    ],
    "refresh" : "",
    "schemaVersion" : 39,
    "tags" : [],
    "templating" : {
      "list" : []
    },
    "time" : {
      "from" : "now-6h",
      "to" : "now"
    },
    "timepicker" : {},
    "timezone" : "",
    "title" : "Monitor Dashboard",
    "uid" : "e115a275-682c-4ec2-8482-55e552c2c3a0",
    "version" : 2,
    "weekStart" : ""
  },
    {
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
      "panels": [
        {
          "datasource": {
            "type": "loki",
            "uid": grafana_data_source.loki.uid
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
                "barWidthFactor": 0.6,
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
          "id": 4,
          "options": {
            "legend": {
              "calcs": [],
              "displayMode": "list",
              "placement": "bottom",
              "showLegend": true
            },
            "tooltip": {
              "hideZeros": false,
              "mode": "single",
              "sort": "none"
            }
          },
          "pluginVersion": "11.5.2",
          "targets": [
            {
              "datasource": {
                "type": "loki",
                "uid": grafana_data_source.loki.uid
              },
              "direction": "backward",
              "editorMode": "builder",
              "expr": "sum(rate({service_name=\"itu-minitwit_loki\"} |= `` [$__auto]))",
              "queryType": "range",
              "refId": "A"
            }
          ],
          "title": "Loki Chart",
          "type": "timeseries"
        },
        {
          "datasource": {
            "type": "loki",
            "uid": grafana_data_source.loki.uid
          },
          "fieldConfig": {
            "defaults": {},
            "overrides": []
          },
          "gridPos": {
            "h": 8,
            "w": 12,
            "x": 12,
            "y": 0
          },
          "id": 3,
          "options": {
            "dedupStrategy": "none",
            "enableInfiniteScrolling": false,
            "enableLogDetails": true,
            "prettifyLogMessage": true,
            "showCommonLabels": false,
            "showLabels": true,
            "showTime": true,
            "sortOrder": "Descending",
            "wrapLogMessage": true
          },
          "pluginVersion": "11.5.2",
          "targets": [
            {
              "datasource": {
                "type": "loki",
                "uid": grafana_data_source.loki.uid
              },
              "direction": "backward",
              "editorMode": "builder",
              "expr": "{service_name=\"itu-minitwit_loki\"} |= ``",
              "queryType": "range",
              "refId": "A"
            }
          ],
          "title": "Loki Logs",
          "type": "logs"
        },
        {
          "datasource": {
            "type": "loki",
            "uid": grafana_data_source.loki.uid
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
                "barWidthFactor": 0.6,
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
            "y": 8
          },
          "id": 1,
          "options": {
            "legend": {
              "calcs": [],
              "displayMode": "list",
              "placement": "bottom",
              "showLegend": true
            },
            "tooltip": {
              "hideZeros": false,
              "mode": "single",
              "sort": "none"
            }
          },
          "pluginVersion": "11.5.2",
          "targets": [
            {
              "datasource": {
                "type": "loki",
                "uid": grafana_data_source.loki.uid
              },
              "direction": "backward",
              "editorMode": "builder",
              "expr": "sum(rate({service_name=\"itu-minitwit_app\"} |= `` [$__auto]))",
              "queryType": "range",
              "refId": "A"
            }
          ],
          "title": "App Chart",
          "type": "timeseries"
        },
        {
          "datasource": {
            "type": "loki",
            "uid": grafana_data_source.loki.uid
          },
          "fieldConfig": {
            "defaults": {},
            "overrides": []
          },
          "gridPos": {
            "h": 8,
            "w": 12,
            "x": 12,
            "y": 8
          },
          "id": 2,
          "options": {
            "dedupStrategy": "none",
            "enableInfiniteScrolling": false,
            "enableLogDetails": true,
            "prettifyLogMessage": true,
            "showCommonLabels": false,
            "showLabels": true,
            "showTime": true,
            "sortOrder": "Descending",
            "wrapLogMessage": true
          },
          "pluginVersion": "11.5.2",
          "targets": [
            {
              "datasource": {
                "type": "loki",
                "uid": grafana_data_source.loki.uid
              },
              "direction": "backward",
              "editorMode": "builder",
              "expr": "{service_name=\"itu-minitwit_app\"} |= ``",
              "queryType": "range",
              "refId": "A"
            }
          ],
          "title": "App Logs",
          "type": "logs"
        },
        {
          "datasource": {
            "type": "loki",
            "uid": "eegtlsawuton4f"
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
                "barWidthFactor": 0.6,
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
            "y": 16
          },
          "id": 5,
          "options": {
            "legend": {
              "calcs": [],
              "displayMode": "list",
              "placement": "bottom",
              "showLegend": true
            },
            "tooltip": {
              "hideZeros": false,
              "mode": "single",
              "sort": "none"
            }
          },
          "pluginVersion": "11.5.2",
          "targets": [
            {
              "datasource": {
                "type": "loki",
                "uid": grafana_data_source.loki.uid
              },
              "direction": "backward",
              "editorMode": "builder",
              "expr": "sum(rate({service_name=\"itu-minitwit_prometheus\"} |= `` [$__auto]))",
              "queryType": "range",
              "refId": "A"
            }
          ],
          "title": "Prometheus Chart",
          "type": "timeseries"
        },
        {
          "datasource": {
            "type": "loki",
            "uid": grafana_data_source.loki.uid
          },
          "fieldConfig": {
            "defaults": {},
            "overrides": []
          },
          "gridPos": {
            "h": 8,
            "w": 12,
            "x": 12,
            "y": 16
          },
          "id": 6,
          "options": {
            "dedupStrategy": "none",
            "enableInfiniteScrolling": false,
            "enableLogDetails": true,
            "prettifyLogMessage": true,
            "showCommonLabels": false,
            "showLabels": true,
            "showTime": true,
            "sortOrder": "Descending",
            "wrapLogMessage": true
          },
          "pluginVersion": "11.5.2",
          "targets": [
            {
              "datasource": {
                "type": "loki",
                "uid": grafana_data_source.loki.uid
              },
              "direction": "backward",
              "editorMode": "builder",
              "expr": "{service_name=\"itu-minitwit_prometheus\"} |= ``",
              "queryType": "range",
              "refId": "A"
            }
          ],
          "title": "Prometheus Logs",
          "type": "logs"
        },
        {
          "datasource": {
            "uid": grafana_data_source.loki.uid
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
                "barWidthFactor": 0.6,
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
            "y": 24
          },
          "id": 7,
          "options": {
            "legend": {
              "calcs": [],
              "displayMode": "list",
              "placement": "bottom",
              "showLegend": true
            },
            "tooltip": {
              "hideZeros": false,
              "mode": "single",
              "sort": "none"
            }
          },
          "pluginVersion": "11.5.2",
          "targets": [
            {
              "direction": "backward",
              "editorMode": "builder",
              "expr": "sum(rate({service_name=\"itu-minitwit_alloy\"} |= `` [$__auto]))",
              "queryType": "range",
              "refId": "A"
            }
          ],
          "title": "Alloy Chart",
          "type": "timeseries"
        },
        {
          "datasource": {
            "type": "loki",
            "uid": grafana_data_source.loki.uid
          },
          "fieldConfig": {
            "defaults": {},
            "overrides": []
          },
          "gridPos": {
            "h": 8,
            "w": 12,
            "x": 12,
            "y": 24
          },
          "id": 8,
          "options": {
            "dedupStrategy": "none",
            "enableInfiniteScrolling": false,
            "enableLogDetails": true,
            "prettifyLogMessage": true,
            "showCommonLabels": false,
            "showLabels": true,
            "showTime": true,
            "sortOrder": "Descending",
            "wrapLogMessage": true
          },
          "pluginVersion": "11.5.2",
          "targets": [
            {
              "datasource": {
                "type": "loki",
                "uid": grafana_data_source.loki.uid
              },
              "direction": "backward",
              "editorMode": "builder",
              "expr": "{service_name=\"itu-minitwit_alloy\"} |= ``",
              "queryType": "range",
              "refId": "A"
            }
          ],
          "title": "Alloy Logs",
          "type": "logs"
        }
      ],
      "preload": false,
      "refresh": "",
      "schemaVersion": 40,
      "tags": [],
      "templating": {
        "list": []
      },
      "time": {
        "from": "now-1h",
        "to": "now"
      },
      "timepicker": {},
      "timezone": "browser",
      "title": "Logging dashboard",
      "uid": grafana_data_source.loki.uid,
      "version": 17,
      "weekStart": ""
    })
}


resource "grafana_user" "helge_and_mircea" {
  name     = "Helge & Mircea"
  login    = "helgeandmircea"
  password = var.helge_and_mircea_password
  is_admin = false
  email    = "ropf@itu.dk"
}