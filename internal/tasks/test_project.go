package tasks

import (
	"fmt"
	"time"
	"zenith-agent/internal/engine"

	"github.com/playwright-community/playwright-go"
)

// ExecuteTestProject - Dummy project for testing auto-discovery
func ExecuteTestProject(bm *engine.BrowserManager) error {
	projectName := "Test Project"
	SetProjectStatus(projectName, true)
	defer SetProjectStatus(projectName, false)

	pw, browser, context, err := bm.CreateContext()
	if err != nil {
		GlobalUpdateStats(projectName, false, fmt.Sprintf("Failed to create context: %v", err))
		return err
	}
	defer pw.Stop()
	defer browser.Close()

	page, err := context.NewPage()
	if err != nil {
		GlobalUpdateStats(projectName, false, fmt.Sprintf("Failed to create page: %v", err))
		return err
	}
	defer page.Close()

	fmt.Println("[TEST PROJECT] Starting dummy execution...")

	// Navigate to example.com
	_, err = page.Goto("https://example.com", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateLoad,
		Timeout:   playwright.Float(30000),
	})
	if err != nil {
		errMsg := fmt.Sprintf("Failed to navigate: %v", err)
		GlobalUpdateStats(projectName, false, errMsg)
		return fmt.Errorf(errMsg)
	}

	// Wait a bit
	time.Sleep(2 * time.Second)

	// Get page title
	title, err := page.Title()
	if err != nil {
		errMsg := fmt.Sprintf("Failed to get title: %v", err)
		GlobalUpdateStats(projectName, false, errMsg)
		return fmt.Errorf(errMsg)
	}

	fmt.Printf("[TEST PROJECT] ✓ Successfully loaded page: %s\n", title)
	GlobalUpdateStats(projectName, true, "")
	
	return nil
}
