package tasks

import (
	"encoding/json"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	statsMutex sync.Mutex
	statsFile  = "stats.json"
)

type HistoryEntry struct {
	Message   string `json:"message"`
	Project   string `json:"project"`
	Type      string `json:"type"` // "success" or "failed"
	Timestamp string `json:"timestamp"`
}

type ProjectStats struct {
	SuccessCount int       `json:"success_count"`
	FailedCount  int       `json:"failed_count"`
	FailReasons  []string  `json:"fail_reasons"`
	IsRunning    bool      `json:"is_running"` // New: Running Status
	LastRun      time.Time `json:"last_run"`   // New: Last Execution Time
}

type GlobalStats struct {
	Projects map[string]*ProjectStats `json:"projects"`
	History  []HistoryEntry           `json:"history"` // New: Activity Log
}

// UpdateStats updates the stats for a specific project
func UpdateStats(projectName string, success bool, reason string) {
	statsMutex.Lock()
	defer statsMutex.Unlock()

	var gs GlobalStats
	gs.Projects = make(map[string]*ProjectStats)

	data, err := os.ReadFile(statsFile)
	if err == nil {
		json.Unmarshal(data, &gs)
	}

	if gs.Projects == nil {
		gs.Projects = make(map[string]*ProjectStats)
	}

	if _, ok := gs.Projects[projectName]; !ok {
		gs.Projects[projectName] = &ProjectStats{FailReasons: []string{}}
	}

	p := gs.Projects[projectName]
	
	// Always update LastRun and assume running since we just finished a task cycle
	p.LastRun = time.Now()
	p.IsRunning = true // Helper logic; in real app main loop controls this, but simplistically true here.

	if success {
		p.SuccessCount++
	} else {
		p.FailedCount++
		cleanReason := reason
		if idx := strings.Index(reason, ":"); idx != -1 {
			cleanReason = reason[:idx]
		}
		
		// Check for duplicate reasons
		exists := false
		for _, r := range p.FailReasons {
			if r == cleanReason {
				exists = true
				break
			}
		}
		if !exists {
			p.FailReasons = append(p.FailReasons, cleanReason)
		}
	}

	// Update History Log
	msg := "Task Executed Successfully"
	typeStr := "success"
	if !success {
		msg = "Task Failed: " + reason
		typeStr = "failed"
		// Truncate msg if too long for UI
		if len(msg) > 60 {
			msg = msg[:57] + "..."
		}
	}

	entry := HistoryEntry{
		Message:   msg,
		Project:   projectName,
		Type:      typeStr,
		Timestamp: time.Now().Format("15:04:05"),
	}

	// Prepend to history (newest first)
	gs.History = append([]HistoryEntry{entry}, gs.History...)
	
	// Keep last 50 entries
	if len(gs.History) > 50 {
		gs.History = gs.History[:50]
	}

	newData, _ := json.MarshalIndent(gs, "", "  ")
	os.WriteFile(statsFile, newData, 0644)
}
