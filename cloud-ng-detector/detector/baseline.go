package main

import (
	"fmt"
	"math"
	"time"
)

// Global variables for baseline stats
var (
	EffectiveMean   float64 = 1.0
	EffectiveStdDev float64 = 1.0
	
	// Global traffic stats (for the whole server)
	GlobalMean      float64 = 1.0
	GlobalStdDev    float64 = 1.0

	// Tracking per-second counts
	perSecondCounts     = make(map[int64]float64) // Per-IP logic uses this via main
	globalSecondCounts  = make(map[int64]float64) // Total server traffic
	
	// Hourly Slots: Map of Hour (0-23) to the Mean recorded during that hour
	HourlyMeans = make(map[int]float64)
	
	lastCalcTime int64
)

// RecordRequest tracks per-IP traffic
func RecordRequest(now int64) {
	perSecondCounts[now]++
}

// RecordGlobalRequest tracks total server traffic
func RecordGlobalRequest(now int64) {
	globalSecondCounts[now]++
}

// RecalculateBaseline runs every 60 seconds (Mark's requirement)
func RecalculateBaseline(now int64) {
	if now-lastCalcTime < 60 {
		return
	}
	lastCalcTime = now

	// 1. Calculate Per-IP Baseline (30-minute window)
	EffectiveMean, EffectiveStdDev = calculateStats(perSecondCounts, now)
	
	// 2. Calculate Global Baseline (30-minute window)
	GlobalMean, GlobalStdDev = calculateStats(globalSecondCounts, now)

	// 3. Update Hourly Slot
	currentHour := time.Now().Hour()
	HourlyMeans[currentHour] = EffectiveMean

	fmt.Printf("[BASELINE] Hourly Slot %d updated | Global Mean: %.2f | IP Mean: %.2f\n", 
		currentHour, GlobalMean, EffectiveMean)
}

// Helper to calculate mean and stddev from a count map
func calculateStats(data map[int64]float64, now int64) (float64, float64) {
	var counts []float64
	startTime := now - 1800 // 30 mins window

	for timestamp, count := range data {
		if timestamp < startTime {
			delete(data, timestamp)
		} else {
			counts = append(counts, count)
		}
	}

	n := float64(len(counts))
	if n == 0 { return 1.0, 1.0 }

	sum := 0.0
	for _, v := range counts { sum += v }
	mean := sum / n

	varianceSum := 0.0
	for _, v := range counts {
		varianceSum += math.Pow(v-mean, 2)
	}
	stdDev := math.Sqrt(varianceSum / n)

	if stdDev == 0 { stdDev = 1.0 }
	return mean, stdDev
}

// CheckGlobalAnomaly checks if total server traffic is spiking
func CheckGlobalAnomaly(totalRate float64) (bool, float64) {
	zScore := (totalRate - GlobalMean) / GlobalStdDev
	// Global alerts only fire on Z-Score > 3.0 (no banning)
	return zScore > 3.0, zScore
}
