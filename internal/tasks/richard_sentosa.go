package tasks

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
	"zenith-agent/internal/engine"

	"github.com/playwright-community/playwright-go"
)

const (
	richardURL        = "https://richardsantosa.com/"
	richardSuccessURL = "https://richardsantosa.com/thanks/"
	richardTimeout    = 60000
)

// generateRichardData generates dummy registration data
func generateRichardData() (name, email, phone string) {
	// Indonesian names
	firstNames := []string{
		"Andi", "Budi", "Cahyo", "Deni", "Eko", "Fajar", "Gilang",
		"Hadi", "Indra", "Joko", "Kevin", "Leo", "Made", "Nico",
		"Omar", "Putra", "Reza", "Sandi", "Toni", "Udin",
	}

	lastNames := []string{
		"Pratama", "Santoso", "Nugroho", "Saputra", "Wibowo",
		"Ramadhan", "Kusuma", "Wijaya", "Setiawan", "Gunawan",
		"Firmansyah", "Hermawan", "Kurniawan", "Maulana", "Hidayat",
	}

	rand.Seed(time.Now().UnixNano())

	firstName := firstNames[rand.Intn(len(firstNames))]
	lastName := lastNames[rand.Intn(len(lastNames))]
	name = fmt.Sprintf("%s %s", firstName, lastName)

	// Generate email (lowercase, no spaces)
	emailName := strings.ToLower(strings.ReplaceAll(name, " ", ""))
	timestamp := time.Now().Unix() % 10000
	email = fmt.Sprintf("%s%d@example.com", emailName, timestamp)

	// Generate Indonesian phone number (62xxx format)
	phoneNumber := rand.Intn(900000000) + 100000000
	phone = fmt.Sprintf("628%d", phoneNumber)

	return name, email, phone
}

// waitForFormReady waits for the form to be visible and ready
func waitForRichardForm(page playwright.Page) error {
	fmt.Println("[RICHARD] Waiting for form to load...")

	// Wait for form container
	formSelector := "#lg-mcrr0myt"

	_, err := page.WaitForSelector(formSelector, playwright.PageWaitForSelectorOptions{
		Timeout: playwright.Float(float64(richardTimeout)),
		State:   playwright.WaitForSelectorStateVisible,
	})

	if err != nil {
		return fmt.Errorf("form not found: %v", err)
	}

	// Additional wait for form inputs to be ready
	time.Sleep(2 * time.Second)

	return nil
}

// fillRichardForm fills the registration form
func fillRichardForm(page playwright.Page, name, email, phone string) error {
	fmt.Printf("[RICHARD] Filling form: %s | %s | %s\n", name, email, phone)

	// Selector for form inputs within the lead generation form
	nameInput := "#lg-mcrr0myt input[name='name']"
	emailInput := "#lg-mcrr0myt input[name='email']"
	phoneInput := "#lg-mcrr0myt input[name='phone']"

	// Fill name
	if err := page.Fill(nameInput, name, playwright.PageFillOptions{
		Timeout: playwright.Float(10000),
	}); err != nil {
		return fmt.Errorf("failed to fill name: %v", err)
	}
	time.Sleep(500 * time.Millisecond)

	// Fill email
	if err := page.Fill(emailInput, email, playwright.PageFillOptions{
		Timeout: playwright.Float(10000),
	}); err != nil {
		return fmt.Errorf("failed to fill email: %v", err)
	}
	time.Sleep(500 * time.Millisecond)

	// Fill phone
	if err := page.Fill(phoneInput, phone, playwright.PageFillOptions{
		Timeout: playwright.Float(10000),
	}); err != nil {
		return fmt.Errorf("failed to fill phone: %v", err)
	}
	time.Sleep(500 * time.Millisecond)

	return nil
}

// submitRichardForm submits the form
func submitRichardForm(page playwright.Page) error {
	fmt.Println("[RICHARD] Submitting form...")

	// Find and click submit button
	submitButton := "#lg-mcrr0myt button[type='submit']"

	err := page.Click(submitButton, playwright.PageClickOptions{
		Timeout: playwright.Float(10000),
	})

	if err != nil {
		// Try alternative selectors
		alternativeSelectors := []string{
			"#lg-mcrr0myt .tve-form-button",
			"#lg-mcrr0myt button.tcb-button-link",
			"#lg-mcrr0myt button:has-text('Kirim')",
			"#lg-mcrr0myt button:has-text('Submit')",
		}

		for _, selector := range alternativeSelectors {
			if clickErr := page.Click(selector, playwright.PageClickOptions{
				Timeout: playwright.Float(5000),
			}); clickErr == nil {
				fmt.Printf("[RICHARD] Clicked submit using selector: %s\n", selector)
				err = nil
				break
			}
		}

		if err != nil {
			return fmt.Errorf("failed to click submit button: %v", err)
		}
	}

	return nil
}

// waitForSuccess waits for redirect to thank you page
func waitForRichardSuccess(page playwright.Page) error {
	fmt.Println("[RICHARD] Waiting for success redirect...")

	// Wait for navigation to thanks page
	err := page.WaitForURL(richardSuccessURL, playwright.PageWaitForURLOptions{
		Timeout: playwright.Float(float64(richardTimeout)),
	})

	if err != nil {
		// Check if we're already on the thanks page
		currentURL := page.URL()
		if strings.Contains(currentURL, "/thanks") {
			fmt.Println("[RICHARD] ✅ Already on thanks page")
			return nil
		}
		return fmt.Errorf("redirect to thanks page failed: %v", err)
	}

	fmt.Println("[RICHARD] ✅ Successfully redirected to thanks page")
	return nil
}

// ExecuteRichardSentosa is the main execution function for Richard Sentosa task
func ExecuteRichardSentosa(bm *engine.BrowserManager) error {
	fmt.Println("\n=== RICHARD SENTOSA BOT: Starting Registration ===")

	// Set project status to running
	SetProjectStatus("RichardSentosa", true)
	defer SetProjectStatus("RichardSentosa", false)

	// Generate dummy data
	name, email, phone := generateRichardData()

	// Initialize browser
	pw, browser, context, err := bm.CreateContext()
	if err != nil {
		return fmt.Errorf("failed to create browser context: %v", err)
	}
	defer pw.Stop()
	defer browser.Close()

	page, err := context.NewPage()
	if err != nil {
		return fmt.Errorf("failed to create new page: %v", err)
	}
	defer page.Close()

	// Navigate to site
	fmt.Printf("[RICHARD] Navigating to %s\n", richardURL)
	_, err = page.Goto(richardURL, playwright.PageGotoOptions{
		Timeout:   playwright.Float(float64(richardTimeout)),
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})
	if err != nil {
		return fmt.Errorf("failed to navigate: %v", err)
	}

	// Wait for form to be ready
	if err := waitForRichardForm(page); err != nil {
		return err
	}

	// Fill form
	if err := fillRichardForm(page, name, email, phone); err != nil {
		return err
	}

	// Submit form
	if err := submitRichardForm(page); err != nil {
		return err
	}

	// Wait for success
	if err := waitForRichardSuccess(page); err != nil {
		return err
	}

	fmt.Println("[RICHARD] ✅ Registration completed successfully!")
	return nil
}
