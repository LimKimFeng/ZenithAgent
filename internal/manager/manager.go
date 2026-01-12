package manager

import (
	"encoding/json"
	"os"
	"time"
)

type AppState struct {
	IsRunning     bool      `json:"is_running"`
	RotatorActive bool      `json:"rotator_active"`
	LastPid       int       `json:"last_pid"`
	StartTime     time.Time `json:"start_time"`
}

const stateFile = "state.json"

func GetState() AppState {
	var state AppState
	data, err := os.ReadFile(stateFile)
	if err == nil {
		json.Unmarshal(data, &state)
	}
	return state
}

func UpdateState(isRunning, rotatorActive bool) {
	state := AppState{
		IsRunning:     isRunning,
		RotatorActive: rotatorActive,
		LastPid:       os.Getpid(),
		StartTime:     time.Now(),
	}
	data, _ := json.MarshalIndent(state, "", "  ")
	os.WriteFile(stateFile, data, 0644)
}
