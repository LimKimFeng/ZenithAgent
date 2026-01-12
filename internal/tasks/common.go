package tasks

import (
	"encoding/json"
	"os"
	"strings"
	"sync"
)

var (
	statsMutex sync.Mutex
	statsFile  = "stats.json"
)

type ProjectStats struct {
	SuccessCount int      `json:"success_count"`
	FailedCount  int      `json:"failed_count"`
	FailReasons  []string `json:"fail_reasons"`
}

type GlobalStats struct {
	Projects map[string]*ProjectStats `json:"projects"`
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

	// Initialize separate project map if nil (for existing file migration or fresh start)
	if gs.Projects == nil {
		gs.Projects = make(map[string]*ProjectStats)
	}

	if _, ok := gs.Projects[projectName]; !ok {
		gs.Projects[projectName] = &ProjectStats{FailReasons: []string{}}
	}

	p := gs.Projects[projectName]
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

	newData, _ := json.MarshalIndent(gs, "", "  ")
	os.WriteFile(statsFile, newData, 0644)
}
