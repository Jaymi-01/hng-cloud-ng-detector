package main

import (
	"fmt"
	"os/exec"
	"time"
)

// StartUnbanner runs continuously in the background
func StartUnbanner() {
	ticker := time.NewTicker(10 * time.Second) // Check every 10 seconds
	
	for range ticker.C {
		now := time.Now()
		
		for ip, record := range ipRecords {
			// If they are banned, AND they aren't permanently banned (Level 4+)
			if record.IsBanned && record.Level <= 3 {
				
				// Has their time expired?
				if now.After(record.UnbanTime) {
					UnbanIP(ip)
				}
			}
		}
	}
}

// UnbanIP executes the release
func UnbanIP(ip string) {
	fmt.Printf("[🕊️ FORGIVER] Releasing IP %s from timeout.\n", ip)

	cmd := exec.Command("iptables", "-D", "DOCKER-USER", "-s", ip, "-j", "DROP")
	cmd.Run()

	ipRecords[ip].IsBanned = false

	logEntry := fmt.Sprintf("UNBAN %s | condition: Time Expired | rate: 0 | baseline: %.2f | duration: N/A", ip, EffectiveMean)
	WriteAuditLog(logEntry)

	SendSlackUnbanAlert(ip)
}
