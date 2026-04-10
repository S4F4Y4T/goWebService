# Comprehensive Observability Master Guide (Enterprise V3.2)

## 1. System Philosophy
This stack ensures **Full-Stack Visibility** (from Host Hardware to Application Trace). We prioritize **DORA metrics** (Availability, Latency, Error Rate) through an integrated LGTM stack (Loki, Grafana, Tempo, Mimir/Prometheus).

---

## 2. Deep Dive: Architecture & Components

### 2.1 Logs: High-Performance Ingestion
- **Tool**: **Loki** + **Promtail**.
- **Ingestion**: Promtail watches the `/var/run/docker.sock` to auto-discover containers.
- **Labeling**: Logs are labeled by `service` and `container`. Use `{service="app"}` to filter.
- **Linking**: A Regex-based **Derived Field** in Grafana links JSON `trace_id` fields directly to Tempo traces for 1-click debugging.

### 2.2 Traces: Distributed Context Propagation
- **Tool**: **Tempo** + **OpenTelemetry (OTel)**.
- **SDK**: Instrumented in `pkg/telemetry/otel.go` for HTTP (mux) and DB (GORM) spans.
- **Propagation**: W3C TraceContext headers. Use `telemetry.NewHTTPClient()` for all outgoing service calls.
- **Sampling**: **20% (0.2 ratio)** to balance detail vs. storage cost.

### 2.3 Metrics: Multi-Layered Monitoring
- **Infrastructure**: **Node Exporter** (Host CPU Load, Memory usage, Disk space, and Network utilization).
- **Application**: Go Runtime metrics (GC, Goroutines, Heap) via OTel Collector.
- **Uptime Probing**: **Blackbox Exporter** simulates a real user at `http://app:6969/health` to verify response status and latency.

### 2.4 Security & Gateway (NGINX)
- **Gateway**: Operates at Port **8080**.
- **Auth**: Prometheus and Alertmanager are protected by **Basic Auth** (`admin` / `admin123`).
- **Grafana**: Accessible at `/grafana/` (Admin/Admin).

---

## 3. Service Inventory & Direct Links

| Service | Internal Port | Gateway Access | Description |
| :--- | :--- | :--- | :--- |
| **Go App** | `6969` | `http://localhost:8080/` | Main API Service |
| **Grafana** | `3000` | `localhost:8080/grafana/` | Visualization Portal |
| **Prometheus**| `9090` | `localhost:8080/prometheus/` | Metrics Engine (Auth Required) |
| **Alertmanager**| `9093` | N/A | Alarm Management Pipeline |
| **OTel Collector**| `4318` | N/A | OTLP Ingestion (HTTP) |
| **Loki** | `3100` | N/A | Aggregated Log Lake |
| **Tempo** | `3200` | N/A | Distributed Trace Store |

---

## 4. Maintenance & Operations

### 4.1 Adding New Alerts
- Edit `config/alert.rules.yml`. Add your logic under `infrastructure_alerts`.
- Always include a `runbook_url` for incident guidance.
- Reload: `docker compose restart prometheus`.

### 4.2 Managing SLOs
- **Availability Alert**: Fires if Success Rate < 99% (2m window).
- **Latency Alert**: Fires if 90% of requests > 500ms (5m window).
- **Watchdog (Dead Man's Switch)**: Must ALWAYS be in "Firing" state. If missing, the alerting system is down.

### 4.3 Log Retention & Cleanup
- **Policy**: 7 Days (168h). Managed in `config/loki-config.yml`.

---

## 5. Persistence & Disaster Recovery
- **Volumes**: Data resides in `web_prometheus_data`, `web_grafana_data`, and `postgres_data`.
- **Dashboard Provisioning**: Dashboards are loaded from `./config/dashboards`. **Warning**: UI changes are ephemeral; update JSON files for permanent changes.

---

## 6. Quick Verification Checklist
1. Visit **[Prometheus Targets](http://localhost:8080/prometheus/targets)**: Ensure all targets (node-exporter, blackbox, etc.) are **UP**.
2. Visit **[Prometheus Alerts](http://localhost:8080/prometheus/alerts)**: Ensure **Watchdog** is firing.
3. Visit **[Grafana Explore](http://localhost:8080/grafana/explore)**: Query Loki and click on a Trace ID to verify links.

---

*Compiled by Antigravity Observability Engine.*
