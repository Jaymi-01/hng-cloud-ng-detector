package main

import (
	"fmt"
	"html/template"
	"net/http"
	"runtime"
	"sort"
	"time"
)

// This HTML template uses basic CSS for a clean "security" look
const dashboardHTML = `
<!DOCTYPE html>
<html>
<head>
    <title>HNG Cloud.ng - Security Dashboard</title>
    <meta http-equiv="refresh" content="3">
    <style>
        body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; background: #1a1a1a; color: #eee; margin: 20px; }
        .card { background: #2d2d2d; padding: 20px; border-radius: 8px; margin-bottom: 20px; border-left: 5px solid #007bff; }
        .alert { border-left-color: #dc3545; }
        table { width: 100%; border-collapse: collapse; }
        th, td { text-align: left; padding: 10px; border-bottom: 1px solid #444; }
        .stat-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 20px; }
        .value { font-size: 24px; font-weight: bold; color: #007bff; }
    </style>
</head>
<body>
    <h1>Cloud.ng Live Anomaly Detection</h1>
    
    <div class="stat-grid">
        <div class="card">
            <div>Global Rate</div>
            <div class="value">{{.GlobalRate}} req/s</div>
        </div>
        <div class="card">
            <div>Baseline Mean</div>
            <div class="value">{{printf "%.2f" .Mean}} req/s</div>
        </div>
        <div class="card">
            <div>Std Dev</div>
            <div class="value">{{printf "%.2f" .StdDev}}</div>
        </div>
        <div class="card">
            <div>CPU / Memory</div>
            <div class="value">{{.CPU}} Cores / {{.Mem}} MB</div>
        </div>
    </div>

    <div class="card alert">
        <h3>Banned IPs</h3>
        <table>
            <tr><th>IP Address</th><th>Level</th><th>Status</th><th>Expires In</th></tr>
            {{range $ip, $rec := .Bans}}
            {{if $rec.IsBanned}}
            <tr>
                <td><code>{{$ip}}</code></td>
                <td>{{$rec.Level}}</td>
                <td>BANNED</td>
                <td>{{$rec.UnbanTime.Sub $.Now | printf "%.0f"}}s</td>
            </tr>
            {{end}}
            {{end}}
        </table>
    </div>

    <div class="card">
        <h3>Top 10 Source IPs (Last 60s)</h3>
        <table>
            <tr><th>IP Address</th><th>Requests</th></tr>
            {{range .TopIPs}}
            <tr><td><code>{{.IP}}</code></td><td>{{.Count}}</td></tr>
            {{end}}
        </table>
    </div>
</body>
</html>
`

type DisplayStats struct {
	GlobalRate float64
	Mean       float64
	StdDev     float64
	CPU        int
	Mem        uint64
	Bans       map[string]*BanRecord
	TopIPs     []IPCount
	Now        time.Time
}

type IPCount struct {
	IP    string
	Count int
}

func StartDashboard() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Calculate Top IPs
		var top []IPCount
		for ip, ts := range ipWindows {
			top = append(top, IPCount{ip, len(ts)})
		}
		sort.Slice(top, func(i, j int) bool { return top[i].Count > top[j].Count })
		if len(top) > 10 { top = top[:10] }

		// System Metrics
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		stats := DisplayStats{
			GlobalRate: float64(len(perSecondCounts)), // Rough global count
			Mean:       EffectiveMean,
			StdDev:     EffectiveStdDev,
			CPU:        runtime.NumCPU(),
			Mem:        m.Alloc / 1024 / 1024,
			Bans:       ipRecords,
			TopIPs:     top,
			Now:        time.Now(),
		}

		tmpl, _ := template.New("dash").Parse(dashboardHTML)
		tmpl.Execute(w, stats)
	})

	fmt.Println("[🖥️  DASHBOARD] Serving live metrics at http://0.0.0.0:8080")
	http.ListenAndServe(":8080", nil)
}
