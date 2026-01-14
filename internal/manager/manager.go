package manager

import (
	"encoding/json"
	"fmt"
	"os"
	"syscall"
	"time"
)

type AppState struct {
	IsRunning     bool      `json:"is_running"`
	RotatorActive bool      `json:"rotator_active"`
	LastPid       int       `json:"last_pid"`
	StartTime     time.Time `json:"start_time"`
}

const (
	stateFile = "state.json"
	lockFile  = ".zenith.lock"
)

var lockFd *os.File

// AcquireLock attempts to acquire an exclusive lock for the application
// Returns error if another instance is already running
func AcquireLock() error {
	var err error
	lockFd, err = os.OpenFile(lockFile, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("failed to open lock file: %w", err)
	}

	// Try to acquire exclusive lock (non-blocking)
	err = syscall.Flock(int(lockFd.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		lockFd.Close()
		return fmt.Errorf("❌ Another instance of ZenithAgent is already running. Please stop it first or wait for it to complete")
	}

	// Write current PID to lock file
	lockFd.Truncate(0)
	lockFd.Seek(0, 0)
	lockFd.WriteString(fmt.Sprintf("%d\n", os.Getpid()))
	lockFd.Sync()

	return nil
}

// ReleaseLock releases the file lock
func ReleaseLock() {
	if lockFd != nil {
		syscall.Flock(int(lockFd.Fd()), syscall.LOCK_UN)
		lockFd.Close()
		os.Remove(lockFile)
	}
}

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
