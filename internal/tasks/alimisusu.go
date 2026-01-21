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
	alimisusuURL     = "https://www.alimisusu28.asia/idn3475"
	alimisusuTimeout = 60000
)

// Checkbox options for combo selection
var alimisusuComboOptions = []string{
	"Harga hanya 322.000RP/1 kotak jika membeli kombo 3 + 1",
	"Harga hanya 265.000RP/1 kotak jika membeli kombo 4 + 2",
	"Harga hanya 245.000RP/1 kotak jika membeli kombo 5 + 3",
	"Harga hanya 230.000RP/1 kotak jika membeli kombo 6 + 4",
}

// generateAlimisusuData generates dummy order data
func generateAlimisusuData() (name, phone, address, selectedCombo string) {
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

	streets := []string{
		"Jl. Merdeka", "Jl. Sudirman", "Jl. Thamrin", "Jl. Gatot Subroto",
		"Jl. Ahmad Yani", "Jl. Diponegoro", "Jl. Imam Bonjol",
		"Jl. Hayam Wuruk", "Jl. Gajah Mada", "Jl. Pemuda",
	}

	cities := []string{
		"Jakarta", "Bandung", "Surabaya", "Semarang", "Yogyakarta",
		"Malang", "Solo", "Medan", "Palembang", "Makassar",
	}

	provinces := []string{
		"DKI Jakarta", "Jawa Barat", "Jawa Timur", "Jawa Tengah",
		"DI Yogyakarta", "Sumatera Utara", "Sumatera Selatan", "Sulawesi Selatan",
	}

	rand.Seed(time.Now().UnixNano())

	// Generate name
	firstName := firstNames[rand.Intn(len(firstNames))]
	lastName := lastNames[rand.Intn(len(lastNames))]
	name = fmt.Sprintf("%s %s", firstName, lastName)

	// Generate phone (Indonesian format)
	phoneNumber := rand.Intn(900000000) + 100000000
	phone = fmt.Sprintf("628%d", phoneNumber)

	// Generate full address
	street := streets[rand.Intn(len(streets))]
	houseNum := rand.Intn(200) + 1
	rtNum := rand.Intn(20) + 1
	rwNum := rand.Intn(15) + 1
	city := cities[rand.Intn(len(cities))]
	province := provinces[rand.Intn(len(provinces))]

	address = fmt.Sprintf("%s No. %d RT %02d/RW %02d, Kel. %s, Kec. %s, %s, %s",
		street, houseNum, rtNum, rwNum,
		firstName+"an", lastName+"an", city, province)

	// Random select combo option
	selectedCombo = alimisusuComboOptions[rand.Intn(len(alimisusuComboOptions))]

	return name, phone, address, selectedCombo
}

// fillAlimisusuForm fills the order form
func fillAlimisusuForm(page playwright.Page, name, phone, address, selectedCombo string) error {
	fmt.Println("[ALIMISUSU] Filling form...")

	// Wait for form to be ready
	time.Sleep(2 * time.Second)

	// Fill name
	nameInput := page.Locator("input[name='name']").First()
	if err := nameInput.Fill(name); err != nil {
		return fmt.Errorf("failed to fill name: %v", err)
	}
	time.Sleep(500 * time.Millisecond)

	fmt.Printf("[ALIMISUSU] Filled name: %s\n", name)

	// Fill phone
	phoneInput := page.Locator("input[name='phone']").First()
	if err := phoneInput.Fill(phone); err != nil {
		return fmt.Errorf("failed to fill phone: %v", err)
	}
	time.Sleep(500 * time.Millisecond)

	fmt.Printf("[ALIMISUSU] Filled phone: %s\n", phone)

	// Fill address
	addressInput := page.Locator("input[name='address']").First()
	if err := addressInput.Fill(address); err != nil {
		return fmt.Errorf("failed to fill address: %v", err)
	}
	time.Sleep(500 * time.Millisecond)

	fmt.Printf("[ALIMISUSU] Filled address: %s\n", address)

	// Select random checkbox
	fmt.Printf("[ALIMISUSU] Selecting combo: %s\n", selectedCombo)

	// Find checkbox by value and click it
	checkboxSelector := fmt.Sprintf("input[name='form_item55'][value='%s']", selectedCombo)
	checkbox := page.Locator(checkboxSelector).First()

	if err := checkbox.Check(playwright.LocatorCheckOptions{
		Force: playwright.Bool(true),
	}); err != nil {
		return fmt.Errorf("failed to select checkbox: %v", err)
	}

	time.Sleep(500 * time.Millisecond)
	fmt.Println("[ALIMISUSU] Checkbox selected successfully")

	return nil
}

// submitAlimisusuForm submits the order
func submitAlimisusuForm(page playwright.Page) error {
	fmt.Println("[ALIMISUSU] Submitting order...")

	// Try to click the first visible submit button
	// There are multiple buttons with ID BUTTON5 and BUTTON6
	submitSelectors := []string{
		"#BUTTON5",
		"#BUTTON6",
		"button[type='submit']",
		".ladi-button",
	}

	for _, selector := range submitSelectors {
		button := page.Locator(selector).First()

		// Check if button is visible
		isVisible, err := button.IsVisible()
		if err == nil && isVisible {
			if err := button.Click(); err == nil {
				fmt.Printf("[ALIMISUSU] Clicked submit button: %s\n", selector)
				return nil
			}
		}
	}

	return fmt.Errorf("failed to find and click submit button")
}

// waitForAlimisusuSuccessPopup waits for success popup alert
func waitForAlimisusuSuccessPopup(page playwright.Page) error {
	fmt.Println("[ALIMISUSU] Waiting for success popup...")

	startTime := time.Now()
	timeout := 30 * time.Second

	for time.Since(startTime) < timeout {
		// Check for popup with success message
		popup := page.Locator(".ladipage-message-box")

		isVisible, err := popup.IsVisible()
		if err == nil && isVisible {
			// Get the message text
			messageText := popup.Locator(".ladipage-message-text")
			text, err := messageText.InnerText()

			if err == nil && strings.Contains(text, "Terima kasih atas pembelian Anda!") {
				fmt.Println("[ALIMISUSU] ✅ Success popup appeared!")
				fmt.Printf("[ALIMISUSU] Message: %s\n", text)
				return nil
			}
		}

		time.Sleep(500 * time.Millisecond)
	}

	return fmt.Errorf("timeout waiting for success popup")
}

// ExecuteAlimisusu is the main execution function
func ExecuteAlimisusu(bm *engine.BrowserManager) error {
	fmt.Println("\n=== ALIMISUSU BOT: Starting Order ===")

	// Set project status to running
	SetProjectStatus("Alimisusu", true)
	defer SetProjectStatus("Alimisusu", false)

	// Generate dummy data
	name, phone, address, selectedCombo := generateAlimisusuData()

	fmt.Printf("[ALIMISUSU] Generated Data:\n")
	fmt.Printf("  Name: %s\n", name)
	fmt.Printf("  Phone: %s\n", phone)
	fmt.Printf("  Address: %s\n", address)
	fmt.Printf("  Selected Combo: %s\n", selectedCombo)

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
	fmt.Printf("[ALIMISUSU] Navigating to %s\n", alimisusuURL)
	_, err = page.Goto(alimisusuURL, playwright.PageGotoOptions{
		Timeout:   playwright.Float(float64(alimisusuTimeout)),
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})
	if err != nil {
		return fmt.Errorf("failed to navigate: %v", err)
	}

	// Fill form
	if err := fillAlimisusuForm(page, name, phone, address, selectedCombo); err != nil {
		return err
	}

	// Submit form
	if err := submitAlimisusuForm(page); err != nil {
		return err
	}

	// Wait for success popup
	if err := waitForAlimisusuSuccessPopup(page); err != nil {
		return err
	}

	fmt.Println("[ALIMISUSU] ✅ Order completed successfully!")
	return nil
}
