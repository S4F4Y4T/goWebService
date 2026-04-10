# Observability Stack Documentation

## Overview
This document describes the comprehensive observability stack integrated into our Go Web Service. The stack provides real-time visibility into the system's performance, health, and behavior through a combined paradigm of Logs, Traces, Metrics, and Alerting.

## Architecture & Components

Our multi-layered observability approach utilizes the following components:

### 1. **Logs (PLG Stack)**
*   **Promtail**: Acts as our localized log shipper. It tails the application log file (`tmp/app.log`) alongside docker container logs and pushes them upstream.
*   **Loki**: A highly efficient log aggregation system designed to store and index logs natively optimized for Grafana visualization.

### 2. **Traces (OpenTelemetry & Tempo)**
*   **OpenTelemetry (OTel) SDK**: Integrated directly into the Go application router and configuration (`pkg/telemetry/otel.go`) to inherently auto-instrument HTTP and GORM database requests.
*   **OTel Collector**: Receives trace and metric data sent from the application, processes it, and exports it to specialized backends.
*   **Tempo**: A high-scale distributed tracing backend where structured trace spans are stored and queried.

### 3. **Metrics (Prometheus Ecosystem)**
*   **Prometheus**: The core robust time-series database and monitoring engine. It recurrently scrapes designated targets to collect numerical data over time.
*   **Node Exporter**: Deployed on the `host` network to gather deep system-level hardware and OS metrics (CPU load, Memory usage, Disk space, and Network utilization).
*   **Blackbox Exporter**: Continuously probes the Go application (`http://host.docker.internal:6969/health`) via HTTP to verify uptime and response status simulating an external user request.
*   **Application Metrics**: Pushed via the OpenTelemetry SDK through the OTel Collector, providing detailed native visibility into Go Garbage Collection, Goroutine counts, and heap usage.

### 4. **Alerting (Alertmanager)**
*   **Alertmanager**: Bound to Prometheus to expertly handle, group, and route raised alerts based on severity.
*   **Alerting Rules**: Pre-configured within `config/alert.rules.yml`. Evaluates conditions like `InstanceDown` or `EndpointDown` dynamically over 1-minute intervals.

### 5. **Visualization (Grafana)**
*   **Grafana**: The centralized 'single pane of glass' dashboard interface. It intrinsically correlates queries across Loki, Tempo, and Prometheus simultaneously.

## Service Ports & Interfaces

| Service | Exposes On | Description / Access |
| :--- | :--- | :--- |
| **Grafana** | `3003` | Main Visualization Dashboard (Auth: `admin` / `admin`). |
| **Prometheus** | `9090` | Direct Time-series query tool and active target health checks interface. |
| **Alertmanager** | `9093` | Alert routing, silencing dashboard, and grouping configuration interface. |
| **Tempo** | `3200` | Tracing endpoint (Backend storage ingestion). |
| **Loki** | `3100` | Log parsing and indexing endpoint. |
| **OTel Collector** | `4318` & `8889` | `4318`: HTTP OTLP receiver. `8889`: Prometheus designated scraping endpoint. |
| **Blackbox Exporter**| `9115` | Internal prober metrics exporter interface. |
| **Node Exporter** | `9100` | Raw Linux system metrics endpoint. |

## Quick Start & Verification

**Check System Health:**
1. Open **[Prometheus Targets](http://localhost:9090/targets)** and verify that `node-exporter`, `otel-collector`, and `blackbox` are in the **`UP`** state.
2. Open **[Prometheus Alerts](http://localhost:9090/alerts)** to observe currently defined monitoring rules.

**Viewing Data in Grafana (http://localhost:3003):**
*   **Logs (Loki)**: Navigate to `Explore`, select the **Loki** data source, and execute queries like `{job="varlogs"}`.
*   **Traces (Tempo)**: Navigate to `Explore`, select the **Tempo** data source, and query by standard `TraceID`. (TraceIDs can be quickly derived from corresponding Loki logs utilizing the embedded `correlation_id`).
*   **Metrics (Prometheus)**: Navigate to `Explore`, select the **Prometheus** data source, and construct queries for `go_gc_duration_seconds` for internal application runtime metrics or `node_cpu_seconds_total` to monitor active OS scaling loads.
