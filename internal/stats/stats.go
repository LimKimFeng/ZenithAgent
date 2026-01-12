package stats

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

const StatsFile = "stats.json"

// DailyStats holds the statistics for the 24-hour cycle
type DailyStats struct {
	ProjectName  string    `json:"project_name"`
	LastReset    time.Time `json:"last_reset"`
	SuccessCount int       `json:"success_count"`
	FailedCount  int       `json:"failed_count"`
	FailReasons  []string  `json:"fail_reasons"`
	mu           sync.Mutex
}

// Global instance
var CurrentStats *DailyStats

func InitStats(projectName string) {
	CurrentStats = &DailyStats{
		ProjectName: projectName,
		LastReset:   time.Now(),
		FailReasons: []string{},
	}
	LoadStats()
}

func LoadStats() {
	if _, err := os.Stat(StatsFile); os.IsNotExist(err) {
		SaveStats()
		return
	}

	data, err := os.ReadFile(StatsFile)
	if err == nil {
		CurrentStats.mu.Lock()
		defer CurrentStats.mu.Unlock()
		json.Unmarshal(data, CurrentStats)
	}
}

func SaveStats() {
	CurrentStats.mu.Lock()
	defer CurrentStats.mu.Unlock()

	data, _ := json.MarshalIndent(CurrentStats, "", "  ")
	os.WriteFile(StatsFile, data, 0644)
}

func RecordSuccess() {
	CurrentStats.mu.Lock()
	CurrentStats.SuccessCount++
	CurrentStats.mu.Unlock()
	SaveStats()
}

func RecordFailure(reason string) {
	CurrentStats.mu.Lock()
	CurrentStats.FailedCount++
	
	// Deduplication
	exists := false
	for _, r := range CurrentStats.FailReasons {
		if r == reason {
			exists = true
			break
		}
	}
	if !exists {
		CurrentStats.FailReasons = append(CurrentStats.FailReasons, reason)
	}
	CurrentStats.mu.Unlock()
	SaveStats()
}

func ResetStats() {
	CurrentStats.mu.Lock()
	CurrentStats.SuccessCount = 0
	CurrentStats.FailedCount = 0
	CurrentStats.FailReasons = []string{}
	CurrentStats.LastReset = time.Now()
	CurrentStats.mu.Unlock()
	SaveStats()
}

func GetReport() string {
	CurrentStats.mu.Lock()
	defer CurrentStats.mu.Unlock()
	
	report := fmt.Sprintf("Report for %s\n", CurrentStats.ProjectName)
	report += fmt.Sprintf("Period: %s to %s\n", CurrentStats.LastReset.Format(time.RFC3339), time.Now().Format(time.RFC3339))
	report += fmt.Sprintf("Success: %d\n", CurrentStats.SuccessCount)
	report += fmt.Sprintf("Failed: %d\n", CurrentStats.FailedCount)
	report += "Failure Reasons:\n"
	if len(CurrentStats.FailReasons) == 0 {
		report += "  None\n"
	} else {
		for _, r := range CurrentStats.FailReasons {
			report += fmt.Sprintf("  - %s\n", r)
		}
	}
	return report
}
