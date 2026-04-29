package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

// Simple struct to read the YAML file (without importing external heavy YAML libraries)
func getSlackWebhook() string {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		fmt.Println("Warning: Could not read config.yaml")
		return ""
	}
	
	// Quick manual parse to avoid extra dependencies for one line
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "slack_webhook_url:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				url := strings.TrimSpace(parts[1])
				// remove quotes if they exist
				url = strings.Trim(url, `"'`)
				return url
			}
		}
	}
	return ""
}

// SendSlackAlert sends a formatted message to Slack
func SendSlackAlert(ip string, condition string, rate float64, baseline float64, duration string) {
	webhookURL := getSlackWebhook()
	if webhookURL == "" {
		fmt.Println("[SLACK] No webhook URL configured, skipping alert.")
		return
	}

	timestamp := time.Now().Format(time.RFC1123)
	
	// Format the message block
	message := fmt.Sprintf("🚨 *ANOMALY DETECTED & IP BANNED* 🚨\n"+
		"• *IP Address:* `%s`\n"+
		"• *Condition:* %s\n"+
		"• *Current Rate:* %.0f req/60s\n"+
		"• *Baseline Mean:* %.2f req/s\n"+
		"• *Ban Duration:* %s\n"+
		"• *Time:* %s", 
		ip, condition, rate, baseline, duration, timestamp)

	// Create the JSON payload Slack expects
	payload := map[string]string{
		"text": message,
	}
	
	jsonData, _ := json.Marshal(payload)

	// Send the POST request
	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("[SLACK ERROR] Failed to send alert: %v\n", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("[SLACK] Alert sent successfully!")
}

// SendSlackUnbanAlert sends a formatted release message to Slack
func SendSlackUnbanAlert(ip string) {
	webhookURL := getSlackWebhook()
	if webhookURL == "" {
		return
	}

	timestamp := time.Now().Format(time.RFC1123)
	message := fmt.Sprintf("✅ *IP UNBANNED* ✅\n"+
		"• *IP Address:* `%s`\n"+
		"• *Condition:* Time Expired\n"+
		"• *Time:* %s", 
		ip, timestamp)

	payload := map[string]string{"text": message}
	jsonData, _ := json.Marshal(payload)

	http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
	fmt.Println("[SLACK] Unban alert sent successfully!")
}
// SendGlobalSlackAlert sends a non-blocking alert when the entire server is under pressure
func SendGlobalSlackAlert(totalRate float64, baseline float64, zScore float64) {
	webhookURL := getSlackWebhook()
	if webhookURL == "" {
		return
	}

	timestamp := time.Now().Format(time.RFC1123)
	message := fmt.Sprintf("⚠️ *GLOBAL TRAFFIC ANOMALY* ⚠️\n"+
		"• *Status:* Server-wide Spike Detected\n"+
		"• *Total Rate:* %.0f req/s\n"+
		"• *Global Baseline:* %.2f req/s\n"+
		"• *Z-Score:* %.2f\n"+
		"• *Action:* No IPs banned (Global threshold only)\n"+
		"• *Time:* %s", 
		totalRate, baseline, zScore, timestamp)

	payload := map[string]string{"text": message}
	jsonData, _ := json.Marshal(payload)

	http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
	fmt.Println("[SLACK] Global alert sent successfully!")
}
