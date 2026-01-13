package tasks

import (
	"encoding/json"
	"os"
	"strings"
	"sync"
	"time"
)

type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Project   string `json:"project"`
	Message   string `json:"message"`
	Type      string `json:"type"` // "success" or "error"
}

type ProjectStats struct {
	IsRunning    bool      `json:"is_running"`
	SuccessCount int       `json:"success_count"`
	FailedCount  int       `json:"failed_count"`
	FailReasons  []string  `json:"fail_reasons"`
	LastRun      time.Time `json:"last_run"`
}

type GlobalStats struct {
	Projects map[string]*ProjectStats `json:"projects"`
	History  []LogEntry               `json:"history"`
}

var (
	statsMutex sync.Mutex
	statsFile  = "stats.json"
)

// GlobalUpdateStats menyimpan data ke format yang dimengerti dashboard
func GlobalUpdateStats(projectName string, success bool, reason string) {
	statsMutex.Lock()
	defer statsMutex.Unlock()

	var gs GlobalStats
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
	p.LastRun = time.Now()
	p.IsRunning = true

	msg := "Eksekusi Berhasil"
	logType := "success"

	if success {
		p.SuccessCount++
	} else {
		p.FailedCount++
		logType = "error"
		msg = "Gagal: " + reason
		
		cleanReason := reason
		if idx := strings.Index(reason, ":"); idx != -1 {
			cleanReason = reason[:idx]
		}
		
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

	// Tambah ke History (Maksimal 50 log)
	newLog := LogEntry{
		Timestamp: time.Now().Format("02 Jan 15:04:05"),
		Project:   projectName,
		Message:   msg,
		Type:      logType,
	}
	gs.History = append([]LogEntry{newLog}, gs.History...)
	if len(gs.History) > 50 {
		gs.History = gs.History[:50]
	}

	newData, _ := json.MarshalIndent(gs, "", "  ")
	os.WriteFile(statsFile, newData, 0644)
}

func SetProjectStatus(projectName string, running bool) {
	statsMutex.Lock()
	defer statsMutex.Unlock()
	
	var gs GlobalStats
	data, err := os.ReadFile(statsFile)
	if err == nil {
		json.Unmarshal(data, &gs)
		if gs.Projects != nil {
			if p, ok := gs.Projects[projectName]; ok {
				p.IsRunning = running
				newData, _ := json.MarshalIndent(gs, "", "  ")
				os.WriteFile(statsFile, newData, 0644)
			}
		}
	}
}
