package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

// BanRecord tracks how many times an IP has attacked, and when their timeout ends
type BanRecord struct {
	Level     int
	UnbanTime time.Time
	IsBanned  bool
}

var ipRecords = make(map[string]*BanRecord)

// getBanDuration calculates the penalty based on the offense level
func getBanDuration(level int) (time.Duration, string) {
	switch level {
	case 1:
		return 10 * time.Minute, "10m"
	case 2:
		return 30 * time.Minute, "30m"
	case 3:
		return 2 * time.Hour, "2h"
	default:
		return 0, "Permanent" // 0 means no unban time
	}
}

func BanIP(ip string, rate float64, zScore float64) {
	record, exists := ipRecords[ip]
	if !exists {
		record = &BanRecord{Level: 0, IsBanned: false}
		ipRecords[ip] = record
	}

	if record.IsBanned {
		return // Already currently banned, don't ban again
	}

	record.Level++
	duration, durationStr := getBanDuration(record.Level)

	if duration > 0 {
		record.UnbanTime = time.Now().Add(duration)
	}
	record.IsBanned = true

	fmt.Printf("[🛡️ BOUNCER] Banning IP %s for %s (Offense Level %d)\n", ip, durationStr, record.Level)

	// Drop them at the Docker front door
	cmd := exec.Command("iptables", "-I", "DOCKER-USER", "-s", ip, "-j", "DROP")
	cmd.Run()

	logEntry := fmt.Sprintf("BAN %s | Z-Score %.2f | rate: %.0f | baseline: %.2f | duration: %s", ip, zScore, rate, EffectiveMean, durationStr)
	WriteAuditLog(logEntry)

	SendSlackAlert(ip, fmt.Sprintf("Z-Score Exceeded (%.2f)", zScore), rate, EffectiveMean, durationStr)
}

func WriteAuditLog(action string) {
	f, err := os.OpenFile("audit.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Failed to open audit log: %v\n", err)
		return
	}
	defer f.Close()

	timestamp := time.Now().Format(time.RFC3339)
	f.WriteString(fmt.Sprintf("[%s] %s\n", timestamp, action))
}
