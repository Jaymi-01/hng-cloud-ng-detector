package main

import (
	"fmt"
	"time"
)

// Shared memory for the sliding window
// These are defined here so all files in 'package main' can see them
var ipWindows = make(map[string][]int64)

func main() {
	fmt.Println("🚀 HNG Cloud.ng Anomaly Detector Active")
	fmt.Println("System initialized. Monitoring /var/log/nginx/hng-access.log...")
	
	// 1. Start the Unban background worker (for the backoff schedule)
	go StartUnbanner()
	
	// 2. Start the Live Metrics Dashboard on :8080
	go StartDashboard()
	
	// 3. Start the Baseline Calculator (runs every 60s as per requirements)
	go func() {
		for {
			RecalculateBaseline(time.Now().Unix())
			time.Sleep(60 * time.Second)
		}
	}()

	// 4. Start tailing logs (this blocks the main thread to keep the app running)
	StartMonitoring()
}

// ProcessTraffic is the bridge between log parsing and security logic
func ProcessTraffic(ip string) {
	currentRate := float64(len(ipWindows[ip]))
	
	// A. Check Global Anomaly (Server-wide spike)
	// globalSecondCounts is managed in baseline.go
	totalRate := float64(len(globalSecondCounts))
	ProcessGlobalStats(totalRate)

	// B. Check Individual IP Anomaly (The "Bouncer" Logic)
	// We wait until we have a baseline (Mean > 1) to avoid false positives at startup
	if EffectiveMean > 1.0 {
		isAnomalous, zScore, condition := CheckIPAnomaly(ip, currentRate)
		
		if isAnomalous {
			// Trigger the block, log, and Slack alert
			BanIP(ip, currentRate, zScore)
			
			// Optional: print to terminal for your Tool-running.png screenshot
			fmt.Printf("[🚨 ALERT] %s blocked. Reason: %s | Z-Score: %.2f\n", ip, condition, zScore)
		}
	}
}
