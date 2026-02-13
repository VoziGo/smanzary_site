# Smanzy Monitoring – User Manual

This guide describes the monitoring stack (Prometheus, Grafana, and exporters) used with the Smanzy Docker Compose setup and how to use it day to day.

---

## 1. Overview

The stack collects and visualizes metrics from:

| Component        | Role |
|-----------------|------|
| **Prometheus**  | Stores time-series metrics and evaluates alerting rules. |
| **Grafana**     | Web UI for dashboards and exploring metrics. |
| **Node Exporter** | Host-level metrics (CPU, memory, disk, network) from the host where containers run. |
| **cAdvisor**    | Per-container CPU, memory, and I/O metrics. |
| **Postgres Exporter** | PostgreSQL metrics (connections, queries, cache, etc.). |
| **Backend (optional)** | Application metrics from the Go backend when `/metrics` is exposed. |

All services run in Docker and communicate over the `smanzy_network` network. Only Prometheus and Grafana are exposed on the host (on localhost).

---

## 2. Ports and URLs

| Service     | Host URL                | Host port | Notes |
|------------|--------------------------|-----------|--------|
| **Grafana**   | http://127.0.0.1:3001   | 3001      | Main UI for dashboards. |
| **Prometheus** | http://127.0.0.1:9091   | 9091      | Query UI and target status. |

Other services (node_exporter, cadvisor, postgres_exporter) are not published; only Prometheus scrapes them inside the network.

---

## 3. Starting and Stopping

- **Start everything (including monitoring):**
  ```bash
  docker compose up -d
  ```

- **Start only the app (no monitoring):**
  ```bash
  docker compose up -d postgres backend frontend thumbnailer
  ```

- **Stop everything:**
  ```bash
  docker compose down
  ```

- **Restart only the monitoring stack:**
  ```bash
  docker compose restart prometheus grafana
  ```

Grafana and Prometheus use named volumes (`grafana_data`, `prometheus_data`), so dashboards and history persist across `docker compose down` and `up`.

---

## 4. Grafana

### 4.1 Login

1. Open **http://127.0.0.1:3001** in your browser.
2. Default credentials (unless overridden by env):
   - **User:** `admin`
   - **Password:** `admin`  
   You are prompted to change the password on first login (recommended).

Override via environment in `docker-compose.yml`:

- `GRAFANA_ADMIN_USER`
- `GRAFANA_ADMIN_PASSWORD`

### 4.2 Prometheus datasource

A Prometheus datasource named **Prometheus** is provisioned automatically and set as default. You do not need to add it manually.

- **Datasource URL (inside Docker):** `http://prometheus:9090`
- If the datasource is missing, check that the file `grafana/provisioning/datasources/prometheus.yml` exists and that the `grafana` service mounts it.

### 4.3 Adding dashboards

1. In Grafana: **☰ → Dashboards → New → Import**.
2. Enter a dashboard ID from Grafana.com (examples below) and load it.
3. Select the **Prometheus** datasource and click **Import**.

**Suggested dashboard IDs:**

| ID    | Name / purpose |
|-------|-----------------|
| **1860** | Node Exporter Full – host CPU, memory, disk, network. |
| **14282** | PostgreSQL – DB connections, transactions, cache. |
| **14283** | PostgreSQL Overview. |
| **193**  | Docker / cAdvisor – container CPU and memory. |
| **3662** | Prometheus 2.0 – Prometheus self-monitoring. |

Use **Datasource: Prometheus** for all of them.

### 4.4 Useful metrics in Grafana

- **Host:** `node_cpu_*`, `node_memory_*`, `node_filesystem_*`, `node_network_*`.
- **Containers:** `container_cpu_*`, `container_memory_*` (from cAdvisor).
- **Postgres:** `pg_*` (from postgres_exporter).
- **Prometheus:** `prometheus_*`, `scrape_*`.

---

## 5. Prometheus

### 5.1 Web UI

- Open **http://127.0.0.1:9091**.
- **Query:** Use the “Graph” or “Table” tab and run PromQL (e.g. `up`, `node_memory_MemAvailable_bytes`).
- **Targets:** **Status → Targets** to see scrape status for each job (prometheus, node, containers, postgres, backend).

### 5.2 Scrape configuration

Scrape config is in **`prometheus.yml`** at the project root, mounted into the container. Main jobs:

- **prometheus** – self-metrics.
- **node** – node_exporter (host).
- **containers** – cAdvisor.
- **postgres** – postgres_exporter.
- **backend** – Go backend (only if it exposes `/metrics`).

After editing `prometheus.yml`, either:

- Restart Prometheus: `docker compose restart prometheus`, or  
- Use lifecycle API (if enabled): `curl -X POST http://127.0.0.1:9091/-/reload`.

### 5.3 Retention

Prometheus is started with `--storage.tsdb.retention.time=15d`. Data older than 15 days is deleted automatically.

---

## 6. Backend application metrics (optional)

The `prometheus.yml` includes a **backend** job scraping `backend:8080`. If your Go backend does **not** expose a `/metrics` endpoint, this target will show as down or return 404. That is expected and does not affect other jobs.

To enable backend metrics:

1. Add a Prometheus client (e.g. `prometheus/client_golang`) to the Go app.
2. Expose an HTTP handler on `/metrics`.
3. No change is needed in `prometheus.yml`; the backend job will then start returning metrics.

---

## 7. Troubleshooting

### Grafana or Prometheus not reachable

- Check that containers are up: `docker compose ps`.
- Check ports: `ss -tlnp | grep -E '3001|9091'` (or equivalent on your OS).

### “Target down” or no data in Grafana

- **Status → Targets** in Prometheus: see which job is failing.
- Ensure dependent services are healthy: `docker compose ps` (postgres, postgres_exporter, cadvisor, node_exporter).
- Restart the failing service, then Prometheus:  
  `docker compose restart <service> prometheus`.

### Grafana “datasource not found”

- Confirm the provisioning file exists:  
  `grafana/provisioning/datasources/prometheus.yml`  
  and that the Grafana service mounts it (see `docker-compose.yml`).
- Restart Grafana: `docker compose restart grafana`.

### cAdvisor permission errors (e.g. on some hosts)

If cAdvisor logs show permission errors reading host paths, you may need to run it with more privileges. In `docker-compose.yml`, under the `cadvisor` service, you can try adding:

```yaml
privileged: true
```

Use only if required and only in environments where you accept the security impact.

### High disk usage by Prometheus

- Reduce retention in the Prometheus `command` in `docker-compose.yml`, e.g.  
  `--storage.tsdb.retention.time=7d`.
- Check volume size: `docker system df -v` and inspect `prometheus_data`.

---

## 8. Improvements applied in `docker-compose.yml`

- **Pinned image tags** for Prometheus, Grafana, node-exporter, postgres-exporter, and cAdvisor to avoid surprise changes from `latest`.
- **Startup order:** `depends_on` so postgres_exporter starts after postgres, Prometheus after exporters, Grafana after Prometheus.
- **Grafana:** `GF_SERVER_ROOT_URL`, configurable admin user/password via env, sign-up disabled, and **Prometheus datasource auto-provisioned**.
- **Grafana healthcheck** so the service is marked healthy once the API responds.
- **Prometheus** `--web.enable-lifecycle` for config reload without full restart (optional).
- **Single network** and no unnecessary host port exposure for exporters.

---

## 9. Production notes

- Set strong `GRAFANA_ADMIN_PASSWORD` (and optionally `GRAFANA_ADMIN_USER`) in production; do not rely on defaults.
- Consider putting Grafana behind a reverse proxy (HTTPS, auth) and binding to a non-public interface.
- For production, consider alerting (Alertmanager) and persistent volumes on a dedicated disk for Prometheus.
- The `docker-compose.prod.yml` in this project does not include the monitoring stack; add the same monitoring services and volumes there if you want them in production.

---

**Quick reference**

- Grafana: http://127.0.0.1:3001 (admin / admin by default).
- Prometheus: http://127.0.0.1:9091.
- Config: `prometheus.yml`, `grafana/provisioning/datasources/prometheus.yml`, `docker-compose.yml`.
