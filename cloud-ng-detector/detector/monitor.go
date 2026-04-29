package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

type LogEntry struct {
	SourceIP     string `json:"source_ip"`
	Status       int    `json:"status"`
	Path         string `json:"path"`
	Method       string `json:"method"`
	ResponseSize int    `json:"response_size"`
}

func StartMonitoring() {
	// In a real scenario, we would load the path from config.yaml
	file, err := os.Open("/var/lib/docker/volumes/HNG-nginx-logs/_data/hng-access.log")
	if err != nil {
		fmt.Printf("Critical Error: Cannot open log file: %v\n", err)
		return
	}
	defer file.Close()

	file.Seek(0, io.SeekEnd)
	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				time.Sleep(100 * time.Millisecond)
				continue
			}
			break
		}

		var entry LogEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			continue
		}

		now := time.Now().Unix()
		
		// 1. Record for Baselines
		RecordRequest(now)
		RecordGlobalRequest(now)

		// 2. Track Errors for Error Surge requirement
		if entry.Status >= 400 {
			IPErrorStats[entry.SourceIP]++
		}

		// 3. Process the window
		ipWindows[entry.SourceIP] = append(ipWindows[entry.SourceIP], now)
		
		// 4. Trigger Check
		ProcessTraffic(entry.SourceIP)
	}
}
