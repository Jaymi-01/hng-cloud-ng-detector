package main

import (
	"fmt"
)

// Track global error baseline
var GlobalErrorMean float64 = 0.5 

// IPErrorStats tracks error counts for specific IPs in the sliding window
var IPErrorStats = make(map[string]int)

// CheckIPAnomaly implements the core detection logic
func CheckIPAnomaly(ip string, currentRate float64) (bool, float64, string) {
	// 1. Determine Thresholds (Check for Error Surge)
	zThreshold := 3.0
	multiplier := EffectiveMean * 5.0
	condition := "Standard Z-Score > 3.0"

	// 2. Apply "Tightening" if Error Surge detected
	// If this IP has high errors compared to global average, lower the Z-threshold
	if float64(IPErrorStats[ip]) > (GlobalErrorMean * 3.0) {
		zThreshold = 1.5 // Tighten threshold to catch them faster
		condition = "Error Surge (Threshold Tightened)"
	}

	// 3. Calculate Z-Score
	zScore := (currentRate - EffectiveMean) / EffectiveStdDev

	// 4. Final Decision
	isAnomalous := zScore > zThreshold || currentRate > multiplier
	
	if currentRate > multiplier {
		condition = "Rate > 5x Baseline Mean"
	}

	return isAnomalous, zScore, condition
}

// CheckGlobalAnomaly tracks the server health as a whole
func ProcessGlobalStats(totalRate float64) {
	zScore := (totalRate - GlobalMean) / GlobalStdDev
	
	if zScore > 3.0 {
		fmt.Printf("[⚠️ GLOBAL ALERT] Server-wide traffic spike! Z-Score: %.2f\n", zScore)
		SendGlobalSlackAlert(totalRate, GlobalMean, zScore)
	}
}
