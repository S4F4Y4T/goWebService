# Comprehensive Observability Stack Documentation

## 1. Overview
This document serves as the master guide for the production-grade observability stack built for the Go Web Service. It ensures visibility into system health, performance bottlenecking, and error tracking through Logs, Traces, Metrics, and Alerting.

---

## 2. Architecture & Detailed Components

### 2.1 Logs (Loki / Promtail)
- **Promtail**: The log fetcher. It talks to the **Docker Socket** to automatically discover all containers. It tags logs with `service` and `container` names for easy filtering.
- **Loki**: The storage engine. Unlike other log systems, it only indexes labels, making it extremely fast and lightweight.
- **Derived Fields**: Configured in Grafana to detect `"trace_id"` in JSON logs and provide a deep-link to the exact request trace in Tempo.

### 2.2 Traces (OpenTelemetry & Tempo)
- **OTel SDK**: Instrumented in `pkg/telemetry/otel.go`. It decorates all HTTP routes and GORM queries with spans.
- **OTel Collector**: The middleman at port `4318`. It receives OTLP data and routes it to Tempo (traces) and Prometheus (metrics). 
- **Tempo**: A massive-scale trace storage. It allows you to see "waterfall" charts showing exactly which SQL query or API call caused a slowdown.

### 2.3 Metrics (Prometheus Ecosystem)
- **Prometheus**: The core time-series engine. It scrapes the OTel Collector, Node Exporter, and Blackbox Exporter.
- **Node Exporter**: Accesses host-level metrics (CPU, Memory, Disk) to monitor the local OS environment.
- **Blackbox Exporter**: Continuously probes `http://app:6969/health` to ensure the service is actually responsive to end-users.

### 2.4 Security & Gateway (NGINX)
- **NGINX Gateway**: Acts as the single entry point at port **8080**. It uses a reverse proxy to route traffic and protects sensitive internal tools (Prometheus/Alertmanager) with **Basic Auth**.

---

## 3. Service Ports & Gateway Mapping

| Service | Internal Port | Access via Gateway | Description |
| :--- | :--- | :--- | :--- |
| **Go App** | `6969` | `http://localhost:8080/` | Main application endpoint. |
| **Grafana** | `3000` | `localhost:8080/grafana/` | Visualization (Admin/Admin). |
| **Prometheus**| `9090` | `localhost:8080/prometheus/` | Metrics / DB (Admin/admin123). |
| **OTel Collector**| `4318` | N/A | Receives OTLP data from Go app. |
| **Loki** | `3100` | N/A | Log storage API. |
| **Tempo** | `3200` | N/A | Trace storage API. |

---

## 4. Maintenance & Management Guide

### 4.1 How to add a new Alert:
1. Open `config/alert.rules.yml`.
2. Add a rule (e.g., High Error Rate). 
3. Run `docker compose restart prometheus`.

### 4.2 How to change Log Retention:
1. Open `config/loki-config.yml`.
2. Locate `retention_period: 168h` (7 days). Update as needed.
3. Run `docker compose restart loki`.

### 4.3 Reliability Monitoring (SLOs)
- **Availability**: Triggers if HTTP 5xx errors exceed 1% over 2 minutes.
- **Latency**: Triggers if 90% of requests exceed 500ms.
- **Check Status**: Visit `http://localhost:8080/prometheus/alerts` (Auth required).

---

## 5. Troubleshooting & Persistence

- **"No Data"**: History resets if volumes are missing. We use `web_prometheus_data`, `web_grafana_data`, and `postgres_data` volumes for persistence. Always check if the time range is set to **"Last 15 minutes"**.
- **Port Conflicts**: NGINX is on `8080` to avoid conflicts with local system servers on port `80`.
- **Log-Trace Link**: If links don't appear in Grafana, ensure your Go code is using `slog.JSONHandler` and that the `trace_id` is present in the output.

---

*Document Status: Verified Production Version.*
