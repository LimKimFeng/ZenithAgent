package tasks

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
	"zenith-agent/internal/engine"

	"github.com/playwright-community/playwright-go"
)

// ExecuteAkademiCrypto - Automation for Akademi Crypto registration
func ExecuteAkademiCrypto(bm *engine.BrowserManager) error {
	projectName := "Akademi Crypto"
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

	// Block resource berat untuk hemat VPS
	page.Route("**/*.{png,jpg,jpeg,gif,webp,css}", func(route playwright.Route) {
		route.Abort()
	})

	fmt.Println("[AKADEMI CRYPTO] Starting execution...")

	// Navigate to the website
	url := "https://akademicrypto.com/"
	fmt.Printf("[AKADEMI CRYPTO] Navigating to %s\n", url)
	
	_, err = page.Goto(url, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateLoad,
		Timeout:   playwright.Float(60000),
	})
	if err != nil {
		errMsg := fmt.Sprintf("Failed to navigate: %v", err)
		GlobalUpdateStats(projectName, false, errMsg)
		return fmt.Errorf(errMsg)
	}

	// Wait for page to load
	time.Sleep(3 * time.Second)

	// Generate random user data
	names := []string{"Andi Pratama", "Budi Santoso", "Cahyo Nugroho", "Deni Saputra", "Eko Wibowo", "Fajar Ramadhan", "Gilang Saputra", "Hendra Wijaya", "Irfan Maulana", "Joko Susilo"}
	selectedName := names[rand.Intn(len(names))]
	phone := fmt.Sprintf("628%d", rand.Intn(900000000)+100000000)
	email := fmt.Sprintf("%s%d@example.com", strings.ReplaceAll(strings.ToLower(selectedName), " ", ""), time.Now().Unix()%1000)

	fmt.Printf("[AKADEMI CRYPTO] Generated data - Name: %s, Email: %s, Phone: %s\n", selectedName, email, phone)

	// Wait for form to be visible
	_, err = page.WaitForSelector("#user_first_name1", playwright.PageWaitForSelectorOptions{
		Timeout: playwright.Float(15000),
	})
	if err != nil {
		errMsg := fmt.Sprintf("Form not found: %v", err)
		GlobalUpdateStats(projectName, false, errMsg)
		return fmt.Errorf(errMsg)
	}

	// Type options for human-like behavior
	typeOptions := playwright.PageTypeOptions{Delay: playwright.Float(100)}

	// Fill Nama Lengkap
	fmt.Println("[AKADEMI CRYPTO] Filling name field...")
	if err := page.Type("#user_first_name1", selectedName, typeOptions); err != nil {
		errMsg := fmt.Sprintf("Failed to fill name: %v", err)
		GlobalUpdateStats(projectName, false, errMsg)
		return fmt.Errorf(errMsg)
	}
	time.Sleep(500 * time.Millisecond)

	// Fill Email
	fmt.Println("[AKADEMI CRYPTO] Filling email field...")
	if err := page.Type("#user_email1", email, typeOptions); err != nil {
		errMsg := fmt.Sprintf("Failed to fill email: %v", err)
		GlobalUpdateStats(projectName, false, errMsg)
		return fmt.Errorf(errMsg)
	}
	time.Sleep(500 * time.Millisecond)

	// Fill Phone Number
	fmt.Println("[AKADEMI CRYPTO] Filling phone field...")
	if err := page.Type("#mepr_phone1", phone, typeOptions); err != nil {
		errMsg := fmt.Sprintf("Failed to fill phone: %v", err)
		GlobalUpdateStats(projectName, false, errMsg)
		return fmt.Errorf(errMsg)
	}
	time.Sleep(500 * time.Millisecond)

	// Password is auto-filled by the website's JavaScript to '123'
	// So we don't need to fill it manually

	// Wait before submitting
	time.Sleep(2 * time.Second)

	// Find and click submit button
	submitSelector := "input.mepr-submit[type='submit']"
	fmt.Println("[AKADEMI CRYPTO] Looking for submit button...")
	
	_, err = page.WaitForSelector(submitSelector, playwright.PageWaitForSelectorOptions{
		Timeout: playwright.Float(5000),
	})
	if err != nil {
		errMsg := fmt.Sprintf("Submit button not found: %v", err)
		GlobalUpdateStats(projectName, false, errMsg)
		return fmt.Errorf(errMsg)
	}

	fmt.Println("[AKADEMI CRYPTO] Clicking submit button...")
	if err := page.Click(submitSelector); err != nil {
		errMsg := fmt.Sprintf("Failed to click submit: %v", err)
		GlobalUpdateStats(projectName, false, errMsg)
		return fmt.Errorf(errMsg)
	}

	// Wait for response
	time.Sleep(5 * time.Second)

	// Check for success or error messages
	errorSelectors := []string{
		".mepr-validation-error",
		".cc-error",
		"#login-error",
		".error-message",
	}

	hasError := false
	errorMessage := ""

	for _, selector := range errorSelectors {
		if errorElem, err := page.QuerySelector(selector); err == nil && errorElem != nil {
			if visible, _ := errorElem.IsVisible(); visible {
				if text, err := errorElem.TextContent(); err == nil && text != "" {
					hasError = true
					errorMessage = text
					break
				}
			}
		}
	}

	if hasError {
		errMsg := fmt.Sprintf("Registration failed: %s", errorMessage)
		fmt.Printf("[AKADEMI CRYPTO] %s\n", errMsg)
		GlobalUpdateStats(projectName, false, errMsg)
		return fmt.Errorf(errMsg)
	}

	// If no error, consider it successful
	fmt.Println("[AKADEMI CRYPTO] ✓ Registration completed successfully!")
	GlobalUpdateStats(projectName, true, "")
	
	return nil
}
