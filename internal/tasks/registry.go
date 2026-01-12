package tasks

import (
	"fmt"
	"time"

	"github.com/playwright-community/playwright-go"
)

// Task is the interface that all projects must implement
type Task interface {
	Name() string
	Execute(ctx playwright.BrowserContext) error
}

// Registry holds all available tasks
var Registry = map[string]Task{
	"ProjectA":       &ProjectA{},
	"ProjectB":       &ProjectB{},
	"ProjectC":       &ProjectC{},
}

// ProjectA implementation
type ProjectA struct{}

func (p *ProjectA) Name() string { return "Project A" }

func (p *ProjectA) Execute(ctx playwright.BrowserContext) error {
	fmt.Println("Executing Project A...")
	page, err := ctx.NewPage()
	if err != nil {
		return fmt.Errorf("could not create page: %v", err)
	}
	// Example stealth navigation
	if _, err = page.Goto("https://bot.sannysoft.com/"); err != nil {
		return fmt.Errorf("could not goto: %v", err)
	}
	// Simulate work
	time.Sleep(2 * time.Second)
	fmt.Println("Project A completed successfully.")
	return nil
}

// ProjectB implementation
type ProjectB struct{}

func (p *ProjectB) Name() string { return "Project B" }

func (p *ProjectB) Execute(ctx playwright.BrowserContext) error {
	fmt.Println("Executing Project B...")
	// Add project B logic here
	return nil
}

// ProjectC implementation
type ProjectC struct{}

func (p *ProjectC) Name() string { return "Project C" }

func (p *ProjectC) Execute(ctx playwright.BrowserContext) error {
	fmt.Println("Executing Project C...")
	// Add project C logic here
	return nil
}
