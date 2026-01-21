package tasks

import (
	"fmt"
	"math/rand"
	"strings" // Added missing import
	"time"
	"zenith-agent/internal/engine"

	"github.com/playwright-community/playwright-go"
)

// Smart Reason Generator
func generateSmartReason() string {
	prefix := []string{"Saya ingin ", "Mau sekali ", "Rencana saya ", "Tujuan saya "}
	action := []string{"belajar ", "mendalami ", "memahami cara ", "mencari tahu "}
	topic := []string{"passive income properti ", "strategi pensiun aman ", "jalur lelang perbankan ", "investasi aset masa depan "}
	suffix := []string{"agar keuangan stabil.", "untuk keluarga saya.", "biar ada penghasilan tambahan.", "secara legal dan benar."}

	rand.Seed(time.Now().UnixNano())
	return prefix[rand.Intn(len(prefix))] + action[rand.Intn(len(action))] + topic[rand.Intn(len(topic))] + suffix[rand.Intn(len(suffix))]
}

func ExecuteTernakProperty(bm *engine.BrowserManager) error {
	// Set status running di dashboard
	SetProjectStatus("TernakProperty", true)

	pw, browser, context, err := bm.CreateContext()
	if err != nil {
		return err
	}
	defer pw.Stop()
	defer browser.Close()

	page, err := context.NewPage()
	if err != nil {
		return err
	}
	defer page.Close()

	// Block resource berat untuk hemat VPS
	page.Route("**/*.{png,jpg,jpeg,gif,webp,css,woff,woff2}", func(route playwright.Route) {
		route.Abort()
	})

	// 1. Navigate to new Depok landing page
	_, err = page.Goto("https://ternakproperty.com/depok1", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateLoad,
	})
	if err != nil {
		return err
	}

	// Wait for form to be visible - looking for main form container
	_, err = page.WaitForSelector("#main-form-items-wdJueGRDoG", playwright.PageWaitForSelectorOptions{
		Timeout: playwright.Float(30000),
	})
	if err != nil {
		return fmt.Errorf("timeout waiting for form: %v", err)
	}

	// Data Dummy
	names := []string{"Andi Pratama", "Budi Santoso", "Cahyo Nugroho", "Deni Saputra", "Eko Wibowo", "Fajar Ramadhan", "Gilang Saputra"}
	selectedName := names[rand.Intn(len(names))]
	phone := fmt.Sprintf("628%d", rand.Intn(900000000)+100000000)
	email := fmt.Sprintf("%s%d@example.com", strings.ReplaceAll(strings.ToLower(selectedName), " ", ""), time.Now().Unix()%1000)

	// 2. Fill Form with Human Delay
	typeOptions := playwright.ElementHandleTypeOptions{Delay: playwright.Float(100)}

	// Wait a bit for form to stabilize
	time.Sleep(1 * time.Second)

	// Fill name field (looking for input fields within form)
	nameInput, err := page.QuerySelector("#main-form-items-wdJueGRDoG input[type='text']")
	if err == nil && nameInput != nil {
		nameInput.Type(selectedName, typeOptions)
	}

	// Fill phone - typically second input or has phone-related attributes
	phoneInputs, _ := page.QuerySelectorAll("#main-form-items-wdJueGRDoG input")
	if len(phoneInputs) > 1 {
		phoneInputs[1].Type(phone, typeOptions)
	}

	// Fill email if exists
	if len(phoneInputs) > 2 {
		phoneInputs[2].Type(email, typeOptions)
	}

	// Fill reason/notes (textarea or last text input)
	textarea, _ := page.QuerySelector("#main-form-items-wdJueGRDoG textarea")
	if textarea != nil {
		textarea.Type(generateSmartReason(), typeOptions)
	}

	// 3. Submit - look for submit button
	time.Sleep(500 * time.Millisecond)
	submitBtn, err := page.QuerySelector("#main-form-items-wdJueGRDoG button[type='submit']")
	if submitBtn != nil {
		err = submitBtn.Click()
	} else {
		// Fallback: click any button with submit-like text
		err = page.Click("#main-form-items-wdJueGRDoG button:has-text('Daftar')")
	}

	if err != nil {
		return fmt.Errorf("failed to submit form: %v", err)
	}

	// 4. Success Detection - wait for URL change or success message
	time.Sleep(2 * time.Second)

	// Check if URL changed or success indicator appeared
	currentURL := page.URL()
	if strings.Contains(currentURL, "success") || strings.Contains(currentURL, "terima-kasih") {
		SetProjectStatus("TernakProperty", false)
		return nil
	}

	// Alternative: check for success message on page
	successEl, _ := page.QuerySelector("text=/terima kasih/i")
	if successEl != nil {
		SetProjectStatus("TernakProperty", false)
		return nil
	}

	// Set status ke false setelah selesai
	SetProjectStatus("TernakProperty", false)

	return nil
}

// Local UpdateStats removed; using common.go version
