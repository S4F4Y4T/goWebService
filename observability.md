# Production Observability Guide — Go Web Service
## Enterprise LGTM Stack (Loki · Grafana · Tempo · Mimir/Prometheus)
### Version 4.0 · Last Revised: April 2026

---

## Table of Contents

1. [Philosophy & Scope](#1-philosophy--scope)
2. [Architecture Overview](#2-architecture-overview)
3. [Component Reference](#3-component-reference)
4. [Instrumentation Guide](#4-instrumentation-guide)
5. [Operations Manual](#5-operations-manual)
6. [Data Sources & Dashboards](#6-data-sources--dashboards)
7. [Alerting & Incident Management](#7-alerting--incident-management)
8. [Debugging Runbooks](#8-debugging-runbooks)
9. [Developer Onboarding](#9-developer-onboarding)
10. [Capacity Planning](#10-capacity-planning)
11. [Security & Compliance](#11-security--compliance)

---

## 1. Philosophy & Scope

### 1.1 Observability Pillars

This system is built on the **Three Pillars of Observability** with an additional fourth pillar for large-scale production:

| Pillar | Tool | Purpose |
| :--- | :--- | :--- |
| **Logs** | Loki + Promtail | Structured event records — *what happened* |
| **Metrics** | Prometheus + OTel | Numeric time-series — *how often / how much* |
| **Traces** | Tempo + OTel | Request journeys — *where time was spent* |
| **Uptime Probing** | Blackbox Exporter | Black-box availability — *are we reachable?* |

### 1.2 Signal Strategy

```
HIGH CARDINALITY events  →  Traces  (Tempo)     — Sampled at 20%
AGGREGATED patterns      →  Metrics (Prometheus)  — All requests
HUMAN-READABLE context   →  Logs    (Loki)        — All requests
EXTERNAL perspective     →  Probes  (Blackbox)    — Every 15s
```

### 1.3 SLO Commitments

| SLO | Target | Alert Threshold | Window |
| :--- | :--- | :--- | :--- |
| **Availability** | 99% success rate | < 99% | 5m rolling |
| **Latency (p90)** | < 500ms | > 500ms | 5m rolling |
| **Uptime** | 99.9% | `probe_success == 0` for 1m | continuous |

---

## 2. Architecture Overview

### 2.1 Full Signal Flow

```
┌─────────────────────────────────────────────────────────────────────┐
│                          EXTERNAL TRAFFIC                           │
└──────────────────────────────┬──────────────────────────────────────┘
                               │ :8080
                    ┌──────────▼──────────┐
                    │   NGINX Gateway      │  ← Basic Auth on /prometheus/
                    │   (config/nginx.conf)│    proxy_pass routing
                    └──┬──────────┬────────┘
                       │          │
             ┌─────────▼──┐  ┌────▼──────────┐
             │  Go App    │  │   Grafana      │
             │  :6969     │  │   :3000        │
             └─────┬──────┘  └───────────────┘
                   │ OTLP/HTTP                    ▲ queries
             ┌─────▼──────────────┐               │
             │  OTel Collector    │         ┌──────┴───────┐
             │  :4318 (receive)   │         │  Datasources │
             │  :8889 (metrics)   │         │  ·Prometheus │
             └──────┬─────┬───────┘         │  ·Loki       │
                    │     │                 │  ·Tempo      │
              traces│     │metrics          └──────────────┘
                    │     │
          ┌─────────▼┐  ┌─▼───────────────┐
          │  Tempo   │  │   Prometheus     │
          │  :3200   │  │   :9090          │◄── scrapes ──► Node Exporter
          └──────────┘  └─────────┬────────┘             ► OTel :8889
                                  │                       ► Blackbox
                                  │ fires alerts          ► Postgres Exporter
                              ┌───▼────────────┐
                              │  Alertmanager  │
                              │  :9093         │──► Slack Webhook
                              └────────────────┘

Log Pipeline:
 Docker Containers → Promtail (docker_sd) → Loki :3100
```

### 2.2 Network Layout

All services run on the default Docker Compose network (`web_default`). Inter-service communication uses container hostnames.

```
Gateway (NGINX)  :8080  → externally accessible
App              :6969  → internal; also mapped for local dev
Grafana          :3000  → internal; served at /grafana/ via gateway
Prometheus       :9090  → internal; served at /prometheus/ via gateway (Basic Auth)
Alertmanager     :9093  → internal only
OTel Collector   :4318  → internal only (OTLP HTTP receiver)
OTel Collector   :8889  → internal only (Prometheus scrape endpoint)
Loki             :3100  → internal only
Tempo            :3200  → internal only
Node Exporter    :9100  → host network
Blackbox         :9115  → internal only
Postgres Exp.    :9187  → internal only
```

---

## 3. Component Reference

### 3.1 NGINX — API & Observability Gateway

**Config**: `config/nginx.conf`

| Route | Backend | Auth |
| :--- | :--- | :--- |
| `/` | `app:6969` | None (public API) |
| `/prometheus/` | `prometheus:9090` | HTTP Basic Auth (`.htpasswd`) |
| `/grafana/` | `grafana:3000` | Grafana's own login |

**Basic Auth credentials** are stored in `config/.htpasswd`. To regenerate:
```bash
htpasswd -c config/.htpasswd admin
# Enter new password when prompted
```

---

### 3.2 OTel Collector — Central Telemetry Router

**Config**: `config/otel-collector.yml`

The Collector is the **single ingest point** for all application telemetry. It decouples the app from backend storage.

```
Receivers:  OTLP (gRPC :4317, HTTP :4318)
Processors: batch  ← groups spans/metrics before export
Exporters:
  - otlp/tempo  → Tempo :4317       (traces)
  - prometheus  → exposed at :8889  (metrics)
  - debug       → stdout            (dev diagnostics)
```

**Pipeline Summary:**
```
App → [OTLP HTTP :4318] → batch → Tempo     (trace pipeline)
App → [OTLP HTTP :4318] → batch → Prometheus scrape endpoint (metrics pipeline)
```

---

### 3.3 Prometheus — Metrics Engine

**Config**: `config/prometheus.yml`
**Retention**: 15 days (`--storage.tsdb.retention.time=15d`)
**Scrape Interval**: 15s global

| Job | Target | What It Scrapes |
| :--- | :--- | :--- |
| `otel-collector` | `otel-collector:8889` | App HTTP metrics, Go runtime metrics |
| `node-exporter` | `host.docker.internal:9100` | Host CPU, RAM, disk, network |
| `blackbox` | via relabeling → `blackbox-exporter:9115` | HTTP probe of `/health` |
| `postgres-exporter` | `postgres-exporter:9187` | DB connections, locks, query stats |

---

### 3.4 Loki — Log Aggregation

**Config**: `config/loki-config.yml`

| Setting | Value |
| :--- | :--- |
| **Listen** | HTTP :3100, gRPC :9096 |
| **Storage** | Local filesystem (`/tmp/loki/`) |
| **Retention** | **7 days (168h)** |
| **Compaction** | Every 10m, 2h delete delay |
| **Cache** | Embedded query cache, 100MB |

**Promtail** (`config/promtail-prod-config.yml`) auto-discovers containers via `docker_sd_configs` from `/var/run/docker.sock` and assigns labels:
- `container` — Docker container name
- `service` — Docker Compose service name
- `logstream` — stdout or stderr

---

### 3.5 Tempo — Distributed Trace Store

**Config**: `config/tempo.yml`

| Setting | Value |
| :--- | :--- |
| **Listen** | HTTP :3200 |
| **Receivers** | OTLP HTTP + gRPC |
| **Storage** | Named Docker volume `tempo_data` → `/tmp/tempo/` |
| **Replication** | 1 (single-node) |
| **WAL Block Duration** | 5m |

Trace data is persisted via the `tempo_data` named Docker volume (mounted at `/tmp/tempo`). Traces survive container restarts.

> **Note on full data loss**: `docker compose down -v` will delete the volume and all stored traces. This is expected for a full environment reset.

---

### 3.6 Grafana — Visualization Layer

**Config**: `config/grafana-datasource.yml`, `config/grafana-dashboards.yml`

Pre-configured datasources (auto-provisioned on start):
- **Prometheus** — `http://prometheus:9090`
- **Loki** — with Derived Field regex to extract `trace_id` and link to Tempo
- **Tempo** — `http://tempo:3200`

Pre-provisioned dashboards (from `./config/dashboards/`):
- `go-otel-metrics.json` — HTTP rates, latency histograms, Go runtime

---

### 3.7 Node Exporter — Host Metrics

Runs in `network_mode: host` to access actual host proc/sys filesystems. Scrapes raw hardware metrics: CPU per-core, memory pages, disk I/O, filesystem usage, network packets.

---

### 3.8 Blackbox Exporter — Synthetic Monitoring

**Config**: `config/blackbox.yml`

Probes `http://app:6969/health` every 15s using module `http_2xx`, validating:
- HTTP response is in the `2xx` class
- Response is received within `5s` timeout
- Supports both `HTTP/1.1` and `HTTP/2.0`

---

### 3.9 Postgres Exporter — Database Intelligence

Connects via `DATA_SOURCE_NAME` DSN. Key metrics exposed:
- `pg_stat_activity_count` — active connections per state
- `pg_stat_user_tables_*` — table-level seq scans, index scans
- `pg_database_size_bytes` — per-database size
- `pg_locks_count` — lock contention indicators

---

### 3.10 Alertmanager — Notification Pipeline

**Config**: `config/alertmanager.yml`

```
Prometheus → fires alert → Alertmanager → route match → Slack receiver
```

Grouping strategy:
- Groups by `alertname` + `severity`
- 30s `group_wait` before first notification
- 5m `group_interval` for new alerts in same group
- 4h `repeat_interval` for persistent unresolved alerts

---

## 4. Instrumentation Guide

### 4.1 SDK Initialization (Go)

**File**: `pkg/telemetry/otel.go`

Both tracer and meter providers are initialized at app startup via `main.go`:

```go
// Initialize tracing
shutdownTracer, err := telemetry.InitTracer()
defer shutdownTracer(ctx)

// Initialize metrics (also starts Go runtime metrics)
shutdownMetrics, err := telemetry.InitMetrics()
defer shutdownMetrics(ctx)
```

**OTLP endpoint** is read from env var `OTEL_EXPORTER_OTLP_ENDPOINT` (default: `otel-collector:4318`).

---

### 4.2 Service Identity

The service is registered with Tempo and Prometheus under the name **`go-web-service`** via the OTel resource:

```go
resource.WithAttributes(
    semconv.ServiceName("go-web-service"),
)
```

This name appears as `exported_job="go-web-service"` in Prometheus metrics and as the service name in Tempo traces. **Do not change this string** without updating the alert rules in `config/alert.rules.yml`.

---

### 4.3 HTTP Instrumentation

All HTTP handlers are wrapped with `otelhttp.NewHandler()`. This automatically generates:
- `http_server_duration` histogram (latency per route)
- `http_server_request_size` / `http_server_response_size`
- Injects `trace_id` into span context for log correlation

---

### 4.4 Database Instrumentation

GORM is instrumented via `otelgorm`. Every SQL query generates a child span with:
- Full SQL statement
- DB operation type (SELECT, INSERT, UPDATE, DELETE)
- Table name
- Duration

---

### 4.5 Go Runtime Metrics

Auto-collected via `go.opentelemetry.io/contrib/instrumentation/runtime`:
- `process.runtime.go.goroutines` — goroutine count
- `process.runtime.go.gc.count` — GC frequency
- `process.runtime.go.mem.heap_alloc` — live heap bytes
- `process.runtime.go.mem.gc_pause_ns` — GC pause duration

Exported every **10 seconds** to OTel Collector via OTLP HTTP.

---

### 4.6 Trace Sampling Strategy

```go
sampler := sdktrace.ParentBased(
    sdktrace.TraceIDRatioBased(0.2), // 20% sampling
)
```

**`ParentBased`** respects the upstream sampling decision — if an upstream caller sets `sampled=true`, this service always samples. If no upstream context, sample 20%.

| Environment | Recommended Ratio | Rationale |
| :--- | :--- | :--- |
| Development | `1.0` (100%) | Full visibility during debug |
| Staging | `0.5` (50%) | Balance visibility vs cost |
| Production | `0.2` (20%) | Current setting — cost-effective |
| High Traffic | `0.05` (5%) | Adjust during traffic spikes |

---

### 4.7 Cross-Service Trace Propagation

For outgoing HTTP calls to downstream services, **always use** `telemetry.NewHTTPClient()` instead of the default `http.Client`:

```go
// ✅ Correct — propagates W3C TraceContext headers
client := telemetry.NewHTTPClient()
resp, err := client.Get("http://other-service/api")

// ❌ Wrong — breaks trace continuity
resp, err := http.Get("http://other-service/api")
```

The `NewHTTPClient()` wraps the transport with `otelhttp.NewTransport()`, automatically injecting `traceparent` and `tracestate` headers.

---

### 4.8 Log Correlation

Application logs should emit `trace_id` and `span_id` as structured JSON fields to enable 1-click Loki→Tempo navigation in Grafana:

```json
{
  "level": "error",
  "msg": "database query failed",
  "trace_id": "4bf92f3577b34da6a3ce929d0e0e4736",
  "span_id": "00f067aa0ba902b7",
  "duration_ms": 1203
}
```

Grafana's Loki Derived Field regex extracts `trace_id` and opens the corresponding Tempo trace directly.

---

## 5. Operations Manual

### 5.1 Compose Commands Reference

**Development Stack (`docker-compose.dev.yml`):**
```bash
# Start everything
docker compose -f docker-compose.dev.yml up -d

# Start with fresh build (after code change)
docker compose -f docker-compose.dev.yml up -d --build app

# Stop — preserve all data volumes
docker compose -f docker-compose.dev.yml down

# Stop — DESTROY all volumes (full reset)
docker compose -f docker-compose.dev.yml down -v

# Restart a single service
docker compose -f docker-compose.dev.yml restart prometheus

# Scale (if app supports horizontal scale)
docker compose -f docker-compose.dev.yml up -d --scale app=3
```

**Production Stack (`docker-compose.prod.yml`):**
```bash
docker compose -f docker-compose.prod.yml up -d
docker compose -f docker-compose.prod.yml down
docker compose -f docker-compose.prod.yml ps
```

**Log Inspection:**
```bash
docker compose logs -f app                      # Follow app logs
docker compose logs --tail=200 otel-collector   # Last 200 lines
docker compose logs --since=30m prometheus      # Last 30 minutes
```

---

### 5.2 Stack Health Check — Step-by-Step Procedure

Run these steps **in order**. If a step fails, stop and debug before continuing.

**Step 1 — Container Health**
```bash
docker compose -f docker-compose.dev.yml ps
```
Expected: All containers in `running` state. Any `Exit` or `Restarting` state requires immediate investigation.

**Step 2 — Network Connectivity**
```bash
curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/health
```
Expected output: `200`

**Step 3 — Prometheus Targets**

Visit: `http://localhost:9090/targets`

All targets must show **State: UP**. Check each:
- `otel-collector` — app metrics
- `node-exporter` — host metrics
- `blackbox/0` — synthetic probe
- `postgres-exporter` — DB metrics

**Step 4 — Watchdog Alert**

Visit: `http://localhost:9090/alerts`

The **Watchdog** alert must be in `Firing` state. If absent, the alerting pipeline is broken — investigate Prometheus and Alertmanager connectivity.

**Step 5 — Log Ingestion**

In Grafana → Explore → Loki datasource:
```logql
{service="app"} | json
```
Logs should appear within the last 5 minutes.

**Step 6 — Trace Ingestion**

In Grafana → Explore → Tempo datasource → **Search** tab → Service: `go-web-service`.
At least some spans should appear, assuming the app received traffic.

**Step 7 — Metrics Flow**

In Grafana → Explore → Prometheus:
```promql
http_server_duration_count{exported_job="go-web-service"}
```
Counter should be non-zero and incrementing.

---

### 5.3 Grafana Dashboard Management

**Creating a New Dashboard:**
1. Navigate to **Dashboards → New → New Dashboard**.
2. Add a visualization panel.
3. Select the appropriate datasource (Prometheus / Loki / Tempo).
4. Use the examples below for common queries.

**Saving Permanently (Critical Step!):**

> Dashboards saved via the UI are stored in the Grafana volume. They are **lost** if you run `docker compose down -v`.

To make a dashboard permanent:
```
Dashboard Settings → JSON Model → Copy All → Save to ./config/dashboards/my-dashboard.json
```
The file is auto-provisioned on stack start via `config/grafana-dashboards.yml`.

**Useful PromQL Queries for New Panels:**

```promql
# HTTP Request Rate (per second)
rate(http_server_duration_count{exported_job="go-web-service"}[5m])

# p99 Latency
histogram_quantile(0.99, sum by (le) (rate(http_server_duration_bucket{exported_job="go-web-service"}[5m])))

# Error Rate (5xx only)
rate(http_server_duration_count{exported_job="go-web-service",http_response_status_code=~"5.."}[5m])

# Live Goroutines
process_runtime_go_goroutines{job="otel-collector"}

# Heap In Use
process_runtime_go_mem_heap_inuse_bytes{job="otel-collector"}

# DB Active Connections
pg_stat_activity_count{datname="go_web_service", state="active"}

# Host CPU Usage
100 - (avg by(instance)(rate(node_cpu_seconds_total{mode="idle"}[5m])) * 100)

# Host Memory Usage %
(1 - (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes)) * 100

# Disk Usage %
(1 - (node_filesystem_avail_bytes{mountpoint="/rootfs"} / node_filesystem_size_bytes{mountpoint="/rootfs"})) * 100
```

---

### 5.4 Alert Rule Management

**Adding a New Alert:**
1. Edit `config/alert.rules.yml`.
2. Add rule under an appropriate group (or create a new group).
3. Always include `summary`, `description`, and `runbook_url` annotations.
4. Apply: `docker compose restart prometheus`

**Alert Rule Template:**
```yaml
- alert: MyNewAlert
  expr: <your_promql_expression>
  for: 5m           # Must be continuously true before firing
  labels:
    severity: warning   # warning | critical | none
  annotations:
    summary: "Short human-readable title"
    description: "Detailed description with {{ $value }} for current value"
    runbook_url: "https://your-wiki/runbooks/my-new-alert"
```

---

### 5.5 Volume & Data Management

| Volume | Contains | Risk |
| :--- | :--- | :--- |
| `postgres_data` | All application data | ⛔ Never delete |
| `prometheus_data` | 15 days of metrics | ⚠️ Recoverable from scraping |
| `grafana_data` | Dashboards, users (if not from files) | ⚠️ Back up JSON files |
| (none — `/tmp`) | Loki chunks, Tempo traces | ⚠️ Ephemeral by design |

**Backing Up Prometheus Data:**
```bash
# Snapshot Prometheus data
curl -XPOST http://localhost:9090/api/v1/admin/tsdb/snapshot
# Find snapshot in: docker volume inspect web_prometheus_data
```

---

## 6. Data Sources & Dashboards

### 6.1 Pre-Provisioned Dashboards

| Dashboard | File | Coverage |
| :--- | :--- | :--- |
| **Go OTel Metrics** | `config/dashboards/go-otel-metrics.json` | HTTP rate, latency, goroutines, GC, heap |

### 6.2 Recommended Grafana Dashboard Imports

For production, import these community dashboards via **Dashboards → Import → ID**:

| ID | Name | Purpose |
| :--- | :--- | :--- |
| `1860` | Node Exporter Full | Host hardware deep-dive |
| `9628` | PostgreSQL Database | DB performance |
| `13407` | Loki Logs Overview | Log volume & error rate |
| `15489` | Docker Container Overview | Per-container resource usage |

### 6.3 Grafana Datasource Architecture

```
Grafana Datasource: Prometheus
  └── URL: http://prometheus:9090
  └── Queries: PromQL — all metrics, alerts, SLOs

Grafana Datasource: Loki
  └── URL: http://loki:3100
  └── Queries: LogQL — log search, error counting
  └── Derived Fields:
        regex: "trace_id":"([a-f0-9]+)"
        link:  Tempo datasource → ${__value.raw}

Grafana Datasource: Tempo
  └── URL: http://tempo:3200
  └── Queries: TraceQL, trace search
  └── Service Graph: enabled (via OTel Collector)
```

---

## 7. Alerting & Incident Management

### 7.1 Active Alert Rules

**File**: `config/alert.rules.yml`

| Alert | Expression | Severity | For |
| :--- | :--- | :--- | :--- |
| `InstanceDown` | `up == 0` | critical | 1m |
| `EndpointDown` | `probe_success == 0` | critical | 1m |
| `HighCPUUsage` | CPU > 70% | warning | 5m |
| `HighMemoryUsage` | RAM > 70% | warning | 5m |
| `HighDiskUsage` | Disk > 70% | warning | 5m |
| `AvailabilitySLOViolation` | Success rate < 99% | critical | 2m |
| `LatencySLOViolation` | p90 > 500ms | warning | 5m |
| `Watchdog` | `vector(1)` | none | 1m |

### 7.2 Slack Integration Setup

1. Go to [api.slack.com/apps](https://api.slack.com/apps) → **Create New App** → **From Scratch**.
2. Enable **Incoming Webhooks** → **Add New Webhook to Workspace**.
3. Select your alerts channel → **Allow**.
4. Copy the Webhook URL (format: `https://hooks.slack.com/services/T.../B.../...`).
5. Paste into `config/alertmanager.yml`:
   ```yaml
   global:
     slack_api_url: 'https://hooks.slack.com/services/YOUR/WEBHOOK/URL'
   ```
6. Restart: `docker compose restart alertmanager`

### 7.3 Full Alertmanager Configuration Reference

**File**: `config/alertmanager.yml`

```yaml
global:
  resolve_timeout: 5m
  slack_api_url: 'https://hooks.slack.com/services/YOUR/WEBHOOK/URL'

# Routing tree — evaluated top-down, first match wins
route:
  group_by: ['alertname', 'severity']
  group_wait: 30s        # Wait before sending first notification for a group
  group_interval: 5m     # Time to wait before sending notification for new alerts
  repeat_interval: 4h    # Re-notify if alert is still firing
  receiver: 'slack-notifications'

  routes:
    # Critical alerts go to a dedicated channel with shorter repeat
    - match:
        severity: 'critical'
      receiver: 'slack-critical'
      repeat_interval: 1h

    # Watchdog never notifies (it's a heartbeat)
    - match:
        severity: 'none'
      receiver: 'null-receiver'

receivers:
- name: 'null-receiver'

- name: 'slack-notifications'
  slack_configs:
  - channel: '#alerts'
    send_resolved: true
    icon_emoji: ':robot_face:'
    title: '{{ template "slack.default.title" . }}'
    text: '{{ template "slack.default.text" . }}'

- name: 'slack-critical'
  slack_configs:
  - channel: '#alerts-critical'
    send_resolved: true
    icon_emoji: ':fire:'
    title: '[CRITICAL] {{ .CommonAnnotations.summary }}'
    text: |
      *Alert:* {{ .CommonAnnotations.summary }}
      *Description:* {{ .CommonAnnotations.description }}
      *Runbook:* {{ .CommonAnnotations.runbook_url }}

# Inhibition — suppress lower severity if higher is already firing
inhibit_rules:
  - source_match:
      severity: 'critical'
    target_match:
      severity: 'warning'
    equal: ['alertname', 'instance']
```

### 7.4 Alert Lifecycle Management

**Creating a Silence** (suppress alerts during planned maintenance):

Via Alertmanager UI (`http://localhost:9093`):
1. Click **Silences** → **New Silence**.
2. Add matchers: e.g., `alertname="HighMemoryUsage"`.
3. Set `Starts At` and `Ends At`.
4. Add creator name and comment.
5. Click **Create**.

Via CLI (amtool):
```bash
# Silence for 2 hours
amtool silence add \
  alertname="HighMemoryUsage" \
  --comment="Redis migration in progress" \
  --duration=2h \
  --author="your-name"

# List active silences
amtool silence query

# Expire a silence immediately
amtool silence expire <silence-id>
```

**Understanding Inhibition Rules:**

Inhibition prevents alert spam during cascading failures. Current rule:
- If a `critical` alert fires for an `instance`, all `warning` alerts for the **same** `alertname` + `instance` are suppressed.
- Example: If `InstanceDown{instance="app"}` fires (critical), `HighCPUUsage{instance="app"}` (warning) is silenced — it's irrelevant if the host is down.

---

## 8. Debugging Runbooks

### 8.1 "App is Down or Unhealthy"

```
SYMPTOMS: HTTP 502/504 from gateway, EndpointDown alert firing, users can't reach API

STEP 1 — Identify container state
  $ docker compose ps app
  → If "Exit" or "Restarting": the app has crashed

STEP 2 — Read crash logs
  $ docker compose logs --tail=200 app
  Common causes:
    - "dial tcp postgres:5432: connection refused" → Postgres not ready
    - "panic: runtime error" → Code bug, check stack trace
    - "bind: address already in use" → Port conflict
    - "permission denied" → Volume/file permissions

STEP 3 — Check dependent services
  $ docker compose ps postgres redis
  If Postgres is unhealthy, fix that first — app depends on it.

STEP 4 — Bypass NGINX, probe app directly
  $ curl -v http://localhost:6969/health
  If this returns 200 but gateway returns 5xx → NGINX misconfiguration

STEP 5 — Restart the app
  $ docker compose -f docker-compose.dev.yml up -d app
  Watch logs immediately:
  $ docker compose logs -f app

STEP 6 — Force rebuild (if code/config changed)
  $ docker compose -f docker-compose.dev.yml up -d --build app

STEP 7 — Check Grafana for patterns
  Go to Grafana → Go App Dashboard → Look for:
    - Goroutine count spike (memory leak / deadlock)
    - Heap memory runaway (OOM kill)
    - Error rate spike before crash (upstream trigger)

ESCALATION: If app keeps crashing, collect a core dump or add more structured
logging before the crash point. Check OS-level OOM killer:
  $ sudo dmesg | grep -i "oom-kill"
```

---

### 8.2 "Loki Has No Data / Logs Missing"

```
SYMPTOMS: Grafana Explore → Loki returns no results

STEP 1 — Check Promtail is running
  $ docker compose ps promtail
  $ docker compose logs --tail=100 promtail

  Known error patterns:
    "dial unix /var/run/docker.sock: permission denied"
    → Fix: Ensure promtail runs as root (user: root in compose)

    "connection refused http://loki:3100"
    → Loki is down, proceed to step 3

    "429 Too Many Requests from Loki"
    → Rate limiting, check loki ingestion limits

STEP 2 — Verify Docker socket access
  $ docker compose exec promtail ls /var/run/docker.sock
  Should return the file path without error.

STEP 3 — Check Loki health
  $ curl http://localhost:3100/ready
  Expected: "ready"

  $ curl http://localhost:3100/metrics | grep loki_ingester
  Should show non-zero ingester metrics.

STEP 4 — Check Loki logs
  $ docker compose logs --tail=100 loki
  Look for: "error writing to store", "compaction failed"

STEP 5 — Test broad Loki query in Grafana
  In Explore → Loki: query {job=~".+"}
  If no results at all, ingest pipeline is broken.
  If results exist, your label filter is wrong — check service label.

STEP 6 — Check Promtail config
  Verify config/promtail-prod-config.yml target:
    clients:
      - url: http://loki:3100/loki/api/v1/push
  The container name must be 'loki' (matches Docker Compose service name).

RECOVERY:
  $ docker compose restart loki promtail
  Wait 60s, then run broad query again.
```

---

### 8.3 "Trace Links Not Working / Traces Missing"

```
SYMPTOMS: Clicking trace_id in Loki doesn't open Tempo, or Tempo shows no data

STEP 1 — Verify OTLP endpoint configuration
  Check docker-compose.dev.yml → app service environment:
    OTEL_EXPORTER_OTLP_ENDPOINT=otel-collector:4318
  ⚠️ Must be container hostname, NOT localhost

STEP 2 — Check OTel Collector is receiving spans
  $ docker compose logs --tail=100 otel-collector
  Look for: "TracesExporter" with span counts
  If empty/error → App is not sending spans

STEP 3 — Verify app telemetry initialization
  Check main.go: InitTracer() must be called at startup.
  Check logs for: "failed to create trace exporter" error.

STEP 4 — Check Tempo is accepting data
  $ curl http://localhost:3200/ready
  Expected: "ready"

  $ docker compose logs --tail=100 tempo
  Look for ingestion errors.

STEP 5 — Verify Grafana Tempo datasource
  Grafana → Configuration → Data Sources → Tempo
  → Click "Save & Test" → Should show "Data source connected"

STEP 6 — Check Loki Derived Field configuration
  Grafana → Configuration → Data Sources → Loki
  → Scroll to "Derived Fields"
  → Regex should be: "trace_id":"([a-f0-9]+)"
  → URL should link to the Tempo datasource with ${__value.raw}

STEP 7 — Check sampling rate
  If traffic is low, 20% sampling may miss short test bursts.
  For debugging, temporarily set TraceIDRatioBased(1.0) in otel.go.

RECOVERY:
  $ docker compose restart otel-collector tempo
  Generate some test traffic:
  $ curl http://localhost:8080/api/users
  Wait 30s, then search in Tempo.
```

---

### 8.4 "Prometheus Target is Down"

```
SYMPTOMS: Red target in /prometheus/targets, InstanceDown alert firing

STEP 1 — Identify which target is down from /prometheus/targets
  Check the "Error" column for the specific connection failure message.

STEP 2 — For otel-collector target
  $ docker compose ps otel-collector
  $ docker compose logs otel-collector
  Verify port 8889 is bound: $ docker compose port otel-collector 8889

STEP 3 — For node-exporter target
  Node Exporter runs with network_mode: host.
  $ curl http://localhost:9100/metrics | head -5
  If unreachable, check if node-exporter container is running:
  $ docker compose ps node-exporter

STEP 4 — For postgres-exporter target
  $ docker compose logs postgres-exporter
  Common: wrong DSN. Verify DB_USER / DB_PASSWORD in .env match Postgres.

STEP 5 — For blackbox target
  $ curl "http://localhost:9115/probe?target=http://app:6969/health&module=http_2xx"
  If probe_success=0, the app health endpoint is failing.

STEP 6 — Reload Prometheus config (without restart)
  $ curl -X POST http://localhost:9090/-/reload
```

---

### 8.5 "High Memory / CPU Alert Firing"

```
STEP 1 — Identify the resource consumer
  # Top processes on host
  $ docker stats --no-stream

STEP 2 — For high Go heap memory
  In Grafana → App Metrics:
    process_runtime_go_mem_heap_inuse_bytes
    process_runtime_go_goroutines
  Goroutine leak → function holding goroutines without returning

STEP 3 — For high CPU
  Check if GC is running too frequently:
    rate(process_runtime_go_gc_count[5m])
  High GC with high allocation rate = memory pressure

STEP 4 — Connection pool saturation
  $ curl http://localhost:9187/metrics | grep pg_stat_activity
  If active connections = max_connections, queries are queuing.

STEP 5 — Postgres slow queries
  Connect to DB:
    $ docker compose exec postgres psql -U postgres -d go_web_service
    SELECT query, mean_exec_time, calls FROM pg_stat_statements
    ORDER BY mean_exec_time DESC LIMIT 10;
```

---

## 9. Developer Onboarding

### 9.1 Prerequisites

| Tool | Minimum Version | Install Command (Ubuntu/Debian) | Verify |
| :--- | :--- | :--- | :--- |
| **Docker Engine** | 20.10.0 | `curl -fsSL https://get.docker.com \| sh` | `docker --version` |
| **Docker Compose** | 2.0.0 | Included in Docker Engine v20.10+ | `docker compose version` |
| **Go** | 1.22 | `apt install golang-go` or `brew install go` | `go version` |
| **Make** | any | `apt install make` | `make --version` |

> **Note on Linux**: After installing Docker, add your user to the docker group to run without sudo:
> ```bash
> sudo usermod -aG docker $USER && newgrp docker
> ```

### 9.2 First-Time Local Setup

```bash
# 1. Clone the repository
git clone <repository-url>
cd web

# 2. Set up environment variables
cp .env.example .env
# Review .env — defaults work for local dev, but change passwords in staging/prod

# 3. Launch the full observability stack
docker compose -f docker-compose.dev.yml up -d

# 4. Wait ~45 seconds for all services to initialize
watch docker compose ps

# 5. Verify the stack
curl http://localhost:8080/health
# Expected: {"status":"ok"} or similar

# 6. Access observability tools
# Grafana:    http://localhost:3003     (user/pass from .env → GRAFANA_USER / GRAFANA_PASSWORD)
# Prometheus:  http://localhost:9090
# Alertmanager: http://localhost:9093
```

### 9.3 Key Files to Understand

```
.env                           ← Your local secrets (NOT in git)
.env.example                   ← Template — what keys exist
docker-compose.dev.yml         ← Full local stack definition
docker-compose.prod.yml        ← Production stack (resource-limited)
pkg/telemetry/otel.go          ← OTel SDK initialization
config/
  prometheus.yml               ← Scrape targets
  alert.rules.yml              ← Alert definitions
  alertmanager.yml             ← Notification routing
  loki-config.yml              ← Log retention, storage
  promtail-prod-config.yml     ← Log ship configuration
  tempo.yml                    ← Trace storage
  otel-collector.yml           ← Telemetry routing pipeline
  nginx.conf                   ← API gateway, auth
  grafana-datasource.yml       ← Datasource provisioning
  dashboards/                  ← Auto-provisioned JSON dashboards
```

### 9.4 Day-to-Day Development Workflow

```bash
# Make code changes, then rebuild only the app
docker compose -f docker-compose.dev.yml up -d --build app

# Tail app logs in real time
docker compose logs -f app

# Run a specific Prometheus query to verify your new metric
curl -s 'http://localhost:9090/api/v1/query?query=http_server_duration_count' | jq .

# Query Loki for recent errors from your service
curl -s -G http://localhost:3100/loki/api/v1/query \
  --data-urlencode 'query={service="app"} |= "error"' | jq .

# Check all service health at once
docker compose ps
```

### 9.5 Making Changes to Alert Rules

```bash
# 1. Edit the rules file
vim config/alert.rules.yml

# 2. Validate syntax (dry-run check via Prometheus)
docker compose exec prometheus promtool check rules /etc/prometheus/alert.rules.yml

# 3. Apply by restarting Prometheus
docker compose restart prometheus

# 4. Verify your new alert appears
open http://localhost:9090/alerts
```

---

## 10. Capacity Planning

### 10.1 Current Resource Limits (Production)

| Service | CPU Limit | Memory Limit | Reserved CPU | Reserved RAM |
| :--- | :--- | :--- | :--- | :--- |
| **Go App** | 0.6 cores | 1 GB | 0.3 cores | 512 MB |
| **PostgreSQL** | 1.0 cores | 2 GB | 0.5 cores | 1 GB |
| **Redis** | 0.2 cores | 512 MB | — | — |
| **Loki** | 0.1 cores | 128 MB | — | — |
| **Promtail** | 0.1 cores | 128 MB | — | — |
| **Grafana** | 0.1 cores | 256 MB | — | — |
| **Total** | ~2.1 cores | ~4.1 GB | — | — |

### 10.2 Scaling Triggers

| Metric | Threshold | Action |
| :--- | :--- | :--- |
| CPU > 70% sustained | 5 minutes | Scale `app` horizontally or increase limit |
| Memory > 80% | 5 minutes | Investigate heap leak, or increase limit |
| Postgres connections > 80 | sustained | Increase `max_connections` or add pgBouncer |
| Prometheus TSDB > 80% disk | — | Reduce retention or add storage |
| Loki chunk > 5GB | — | Reduce retention from 7d → 5d |
| p99 latency > 1s | 10 minutes | Profile slow endpoints, add DB indexes |

### 10.3 Storage Growth Estimates

| Data | Growth Rate | 7-Day Total | 15-Day Total |
| :--- | :--- | :--- | :--- |
| Prometheus metrics | ~50 MB/day | ~350 MB | ~750 MB |
| Loki logs (app) | ~100 MB/day | ~700 MB | N/A (7d retention) |
| Tempo traces (20%) | ~200 MB/day | ~1.4 GB | N/A (ephemeral) |
| PostgreSQL WAL | ~100 MB/day | ~700 MB | ~1.5 GB |

---

## 11. Security & Compliance

### 11.1 Credential Management

**Environment File (`.env`):**
```bash
# Database
DB_USER=postgres
DB_PASSWORD=<strong-random-password>
DB_NAME=go_web_service

# Grafana Admin
GRAFANA_USER=admin
GRAFANA_PASSWORD=<strong-random-password>

# Slack Alerting (from config/alertmanager.yml)
# SLACK_WEBHOOK_URL is set directly in alertmanager.yml
```

> **Dev vs Production — Grafana credentials:**
>
> | Environment | Credential source | Default |
> | :--- | :--- | :--- |
> | **Development** | `.env` → `GF_ADMIN_USER` / `GF_ADMIN_PASSWORD` | `admin` / `admin` (`.env.example` default — change it) |
> | **Production** | `.env` → `GRAFANA_USER` / `GRAFANA_PASSWORD` | **Must be set to a strong password** |
>
> The `admin/admin` default only applies if you start the dev stack without editing `.env`. In all other environments, set explicit strong values in `.env` before starting the stack.

**Rules:**
- `.env` is in `.gitignore` — **never commit it**.
- `.env.example` is committed — it contains only key names, never values.
- For CI/CD pipelines, inject secrets via environment variables or Docker Secrets, not files.
- Rotate passwords every 90 days in production.


**NGINX Basic Auth (`.htpasswd`):**
```bash
# Create/update htpasswd file
htpasswd -c config/.htpasswd admin
# Enter strong password at prompt

# Verify it was created
cat config/.htpasswd
```

### 11.2 Port Exposure Policy

| Port | Service | External Exposure | Policy |
| :--- | :--- | :--- | :--- |
| **8080** | NGINX Gateway | ✅ **Allowed** | Only public-facing port |
| `3000/3003` | Grafana | ❌ Block on firewall | Access via `/grafana/` route |
| `9090` | Prometheus | ❌ Block on firewall | Protected via `/prometheus/` + Basic Auth |
| `9093` | Alertmanager | ❌ Block on firewall | Internal only |
| `5432` | PostgreSQL | ❌ **NEVER expose** | Internal + SSH tunnel if needed |
| `6379` | Redis | ❌ **NEVER expose** | Internal only |
| `3100` | Loki | ❌ Block on firewall | Internal only |
| `3200` | Tempo | ❌ Block on firewall | Internal only |
| `4318` | OTel Collector | ❌ Block on firewall | Internal only |
| `9100` | Node Exporter | ❌ Block on firewall | Host-only |

**Recommended UFW rules for production server:**
```bash
sudo ufw default deny incoming
sudo ufw allow 22/tcp    # SSH
sudo ufw allow 8080/tcp  # API Gateway
sudo ufw enable
```

### 11.3 What Goes Through the Gateway (NGINX)

```
Public Internet
      │
      │ :8080
      ▼
  NGINX Gateway
      │
      ├── /           → App (no auth)
      ├── /prometheus/ → Prometheus (HTTP Basic Auth via .htpasswd)
      └── /grafana/   → Grafana (Grafana's own login page)
```

All other services are **not routable** through the gateway. They are internal-only. Use SSH port forwarding for emergency direct access to internal services:

```bash
# Access Alertmanager internally (via SSH tunnel)
ssh -L 9093:localhost:9093 user@your-server
open http://localhost:9093
```

### 11.4 Audit Checklist

Before deploying to production, verify:

- [ ] `.env` is not tracked in git (`git status` shows it as untracked)
- [ ] `config/.htpasswd` has a strong password (not `admin123`)
- [ ] Grafana admin password changed from default
- [ ] UFW rules configured — only port 8080 open
- [ ] Slack webhook is set and tested
- [ ] Watchdog alert is firing in Prometheus
- [ ] All 4 Prometheus targets are `UP`
- [ ] Loki Derived Fields configured for trace correlation
- [ ] Dashboards saved as JSON files in `config/dashboards/`

---

## Quick Reference Card

```
STACK MANAGEMENT
  Start dev:     docker compose -f docker-compose.dev.yml up -d
  Stop dev:      docker compose -f docker-compose.dev.yml down
  Build+start:   docker compose -f docker-compose.dev.yml up -d --build app
  Status:        docker compose ps

ACCESS POINTS
  API:           http://localhost:8080/
  Grafana:       http://localhost:3003     (admin/admin) (see .env → GRAFANA_USER / GRAFANA_PASSWORD)
  Prometheus:    http://localhost:9090
  Alertmanager:  http://localhost:9093
  Loki:          http://localhost:3100/ready
  Tempo:         http://localhost:3200/ready

HEALTH SIGNALS
  All green if:  /prometheus/targets → all UP
                 /prometheus/alerts  → Watchdog firing
                 Grafana Explore → {service="app"} returns logs

CRITICAL FILES
  Alerts:        config/alert.rules.yml
  Routing:       config/alertmanager.yml
  Retention:     config/loki-config.yml  (168h = 7 days)
  Secrets:       .env (never commit)
  Dashboards:    config/dashboards/*.json (auto-provisioned)
```

---

*Compiled by Antigravity Observability Engine · Enterprise V4.0 · April 2026*
