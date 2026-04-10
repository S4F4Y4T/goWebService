# Observability Stack Documentation

## Overview
This document describes the comprehensive observability stack integrated into our Go Web Service. The stack provides real-time visibility into the system's performance, health, and behavior through a combined paradigm of Logs, Traces, Metrics, and Alerting.

## Architecture & Components

Our multi-layered observability approach utilizes the following components:

### 1. **Logs (PLG Stack)**
*   **Docker Ingestion**: Promtail is configured to talk directly to the Docker socket (`/var/run/docker.sock`). It automatically discovers all running containers and scrapes their `stdout` logs.
*   **Loki**: A highly efficient log aggregation system designed to store and index logs natively optimized for Grafana visualization. Logs are labeled by `service` and `container`.
*   **Derived Fields**: Grafana is configured to automatically detect TraceIDs inside JSON logs and provide a one-click link to the corresponding Waterfall trace in Tempo.

### 2. **Traces (OpenTelemetry & Tempo)**
*   **OpenTelemetry (OTel) SDK**: Integrated directly into the Go application router (`pkg/telemetry/otel.go`) to auto-instrument HTTP and GORM database requests.
*   **OTel Collector**: The central hub that receives OTLP data from the app (port 4318) and shovels it to the backends.
*   **Tempo**: Stores and queries structured trace spans, allowing you to see exactly how much time each DB query or function call took during a single request.

### 3. **Metrics (Prometheus Ecosystem)**
*   **Prometheus**: The core time-series database. It is configured with **Persistence**, so metrics are saved to disk in the `web_prometheus_data` volume.
*   **Node Exporter**: Gathers system-level hardware metrics (CPU, RAM, Disk).
*   **Blackbox Exporter**: Probes the application internally (`http://app:6969/health`) to verify status simulating uptime checks.
*   **Application Metrics**: Global Runtime metrics (Memory, Goroutines, GC, DB Pool) are exported via OTel and scraped by Prometheus.

### 4. **Alerting (Slack Integration)**
*   **Alert Rules**: Configured in `config/alert.rules.yml`. Active alerts include:
    *   **High CPU/Memory/Disk**: Triggers when usage exceeds **70%** for more than 5 minutes.
    *   **InstanceDown**: Triggers if a service container stops.
    *   **EndpointDown**: Triggers if the `/health` check fails.
*   **Slack Routing**: Alertmanager is pre-configured to send notifications to Slack. 
    *   *Note: Requires a valid Slack Webhook URL in `config/alertmanager.yml`.*

### 5. **Visualisation (Custom Grafana Dashboards)**
*   **Go Web Service - Complete Metrics**: A custom-built dashboard that maps modern OTel metrics to the classic Go runtime views (Heap, Stack, GC, Objects).
*   **Go Web Service - Health & Uptime**: A dedicated fail-proof dashboard for uptime, response latency, and HTTP status history.

## Service Ports & Interfaces

| Service | Exposes On | Description / Access |
| :--- | :--- | :--- |
| **Grafana** | `3003` | Main Visualization Dashboard (Auth: `admin` / `admin`). |
| **Prometheus** | `9090` | Direct Time-series query tool and active target health checks interface. |
| **Alertmanager** | `9093` | Alert routing and grouping configuration interface. |
| **Tempo** | `3200` | Tracing backend storage. |
| **Loki** | `3100` | Log parsing and indexing endpoint. |
| **OTel Collector** | `4318` & `8889` | `4318`: OTLP HTTP receiver. `8889`: Prometheus scraping endpoint. |
| **Blackbox Exporter**| `9115` | Internal endpoint prober. |

## Persistence Note
All observability data (Dashboards, Metrics, and Postgres Data) is now saved to **Docker Volumes**. You can safely run `docker compose down` and when you bring the stack back up, all your history and graph settings will be exactly where you left them.

## Viewing Data in Grafana (http://localhost:3003)

1.  **Check Health**: Go to the **Health & Uptime** dashboard to see the live status of the service.
2.  **Debug Requests**: Go to `Explore` -> `Loki`, find an error log, and click the **TraceID** button to see the waterfall chart in Tempo.
3.  **Monitor Performance**: Open **Complete Metrics** and set the time range to **"Last 15 minutes"** to see live memory and goroutine data.
