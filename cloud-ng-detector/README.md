# 🛡️ HNG Cloud.ng Anomaly Detector

An intelligent, self-learning anomaly detection system built to protect high-traffic cloud infrastructure. This tool monitors Nginx access logs in real-time, establishes a rolling statistical baseline, and automatically mitigates DDoS and brute-force attacks using `iptables` and a tiered backoff ban schedule.

---

## 📍 Live Project Links

| Service | URL |
|---|---|
| Nextcloud Service | `http://13.53.235.117` |
| Metrics Dashboard | `http://13.53.235.117:8080` |

---

## 🚀 Technical Architecture

The system is composed of several decoupled Go modules that handle the full lifecycle of a request — from ingestion to mitigation:

| Module | Responsibility |
|---|---|
| **Monitor** | Continuously tails and parses JSON-formatted Nginx logs from a shared Docker volume (`HNG-nginx-logs`). |
| **Baseline** | Computes Mean (μ) and Standard Deviation (σ) over a 30-minute rolling window, updated every 60 seconds. Supports "Hourly Slots" for time-of-day traffic awareness. |
| **Detector** | Calculates Z-Scores and monitors for 4xx/5xx error surges to dynamically tighten security thresholds. |
| **Blocker** | Interacts with the Linux kernel via `iptables` on the `DOCKER-USER` chain to drop malicious packets. |
| **Unbanner** | Implements a tiered backoff schedule (10m → 30m → 2h → Permanent) with automated Slack notifications. |
| **Dashboard** | A real-time web UI serving system metrics, global request rates, and active bans. |

---

## 🛠️ Key Design Decisions

### Language Choice: Go (Golang)

Go was selected for its superior performance in networking and systems programming. Unlike Python, Go provides native multi-threading via goroutines, allowing the system to parse logs, calculate complex statistics, and serve a web UI **concurrently** — without the overhead of a Global Interpreter Lock. This ensures the detector itself never becomes a bottleneck during an active attack.

### Sliding Window Logic

Traffic is tracked using a **Deque-based Sliding Window**:

- **Structure:** A map of IP addresses to slices of Unix timestamps — `map[string][]int64`
- **Eviction:** Every time a log line is processed, timestamps older than 60 seconds are evicted for that IP, yielding a precise real-time **Requests Per Minute (RPM)** metric.

### Statistical Baseline & Anomaly Detection

Rather than relying on static thresholds, the system uses a **Dynamic Rolling Baseline**:

- **Window Size:** 30 minutes of per-second traffic data
- **Recalculation Interval:** Every 60 seconds
- **Flagging Conditions:**
  - IP rate exceeds a **Z-Score of 3.0**, or
  - IP rate is **> 5× the baseline mean**
- **Error Surge:** If an IP's 4xx/5xx error rate is **3× the global error average**, the Z-Score threshold is automatically tightened to catch "noisy" attackers faster.

---

## 📦 Setup & Deployment (Fresh VPS)

### 1. System Preparation

```bash
sudo apt update && sudo apt install -y docker.io docker-compose golang-go git
```

### 2. Provision the Stack

Clone the repository and launch the containerized Nextcloud service:

```bash
git clone https://github.com/Jaymi-01/hng-cloud-ng-detector.git
cd hng-cloud-ng-detector
sudo docker-compose up -d
```

### 3. Run the Detector

Configure your Slack Webhook in `detector/config.yaml`, then launch the engine:

```bash
cd detector
sudo go run *.go
```

---

## 🛡️ Audit Log Format

All security events are recorded in a structured `audit.log`:

```
[TIMESTAMP] ACTION ip | condition | rate | baseline | duration
```
