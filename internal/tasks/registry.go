package tasks

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"zenith-agent/internal/engine"

	"github.com/playwright-community/playwright-go"
)

// Task is the interface that all projects must implement
type Task interface {
	Name() string
	Execute(ctx playwright.BrowserContext) error
}

// TaskInfo holds information about a discovered task
type TaskInfo struct {
	Key         string                                  // "TernakProperty"
	DisplayName string                                  // "Ternak Property"
	FileName    string                                  // "ternak_property.go"
	Function    func(*engine.BrowserManager) error      // ExecuteTernakProperty
}

// Registry holds all available tasks (will be populated by DiscoverTasks)
var Registry map[string]TaskInfo

// DiscoverTasks scans internal/tasks directory and builds task registry
func DiscoverTasks() (map[string]TaskInfo, error) {
	tasks := make(map[string]TaskInfo)
	
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %v", err)
	}
	
	// Path to tasks directory
	tasksDir := filepath.Join(cwd, "internal", "tasks")
	
	// Read directory
	files, err := os.ReadDir(tasksDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read tasks directory: %v", err)
	}
	
	// Scan for task files
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		
		fileName := file.Name()
		
		// Skip non-Go files and helper files
		if !strings.HasSuffix(fileName, ".go") {
			continue
		}
		
		// Skip registry, common, and stats files
		if fileName == "registry.go" || fileName == "common.go" || fileName == "stats.go" {
			continue
		}
		
		// Extract project name from filename
		// Example: "ternak_property.go" -> "TernakProperty"
		baseName := strings.TrimSuffix(fileName, ".go")
		
		// Convert snake_case to PascalCase
		parts := strings.Split(baseName, "_")
		var pascalCase string
		for _, part := range parts {
			if len(part) > 0 {
				pascalCase += strings.ToUpper(part[:1]) + part[1:]
			}
		}
		
		// Create display name (Title Case with spaces)
		displayName := strings.Join(parts, " ")
		displayName = strings.Title(displayName)
		
		// Map to execution function
		var execFunc func(*engine.BrowserManager) error
		
		switch pascalCase {
		case "TernakProperty":
			execFunc = ExecuteTernakProperty
		case "ScalpingSr":
			execFunc = ExecuteScalpingSR
		case "AkademiCrypto":
			execFunc = ExecuteAkademiCrypto
		case "TestProject":
			execFunc = ExecuteTestProject
		default:
			// Skip unknown tasks
			fmt.Printf("[REGISTRY] Warning: No executor found for %s (file: %s)\n", pascalCase, fileName)
			continue
		}
		
		// Add to registry
		tasks[pascalCase] = TaskInfo{
			Key:         pascalCase,
			DisplayName: displayName,
			FileName:    fileName,
			Function:    execFunc,
		}
		
		fmt.Printf("[REGISTRY] Discovered task: %s (%s)\n", displayName, fileName)
	}
	
	return tasks, nil
}

// InitRegistry initializes the global Registry variable
func InitRegistry() error {
	discovered, err := DiscoverTasks()
	if err != nil {
		return err
	}
	Registry = discovered
	return nil
}

// GetTaskList returns a sorted list of task keys for display
func GetTaskList() []string {
	var keys []string
	for key := range Registry {
		keys = append(keys, key)
	}
	return keys
}

// GetTask returns a TaskInfo by key
func GetTask(key string) (TaskInfo, bool) {
	task, exists := Registry[key]
	return task, exists
}
